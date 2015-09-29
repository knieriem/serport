package encoding

import (
	"bytes"
	"io"
)

type Encoder interface {
	Encode([]byte) []byte
}

type Decoder interface {
	Decode([]byte) int
}

type StripCR struct {
	wasCR bool
}

func (d *StripCR) Decode(src []byte) int {
	dest := src

	cr := d.wasCR
	iw := 0
	iLast := len(src) - 1
loop:
	for i, c := range src {
		switch c {
		case '\r':
			cr = true
			if i == iLast {
				break loop
			}
		default:
			cr = false
		case '\n':
			if cr && iw > 0 {
				iw--
				cr = false
			}
		}
		if iw != i {
			dest[iw] = c
		}
		iw++
	}
	d.wasCR = cr
	return iw
}

type TermInput struct {
	useStripCR bool
	StripCR
}

func (d *TermInput) Decode(src []byte) int {
	dest := src

	if d.useStripCR {
		return d.StripCR.Decode(src)
	}
	if bytes.IndexByte(src, '\n') != -1 {
		d.useStripCR = true
		return d.StripCR.Decode(src)
	}
	for i, c := range src {
		if c == '\r' {
			dest[i] = '\n'
		}
	}
	return len(dest)
}

type InsertCR struct {
	buf []byte
}

func (e *InsertCR) Encode(src []byte) []byte {
	var dest []byte
	push := func(b byte) {
		if dest != nil {
			dest = append(dest, b)
		}
	}
	for i, c := range src {
		if c == '\n' {
			if dest == nil {
				if rsrvCap := len(src) * 2; cap(e.buf) < rsrvCap {
					e.buf = make([]byte, rsrvCap)
				}
				dest = e.buf[:i]
				copy(dest, src)
			}
			push('\r')
		}
		push(c)
	}
	if dest != nil {
		e.buf = dest
	} else {
		dest = src
	}
	return dest
}

type Wrapper struct {
	r io.Reader
	w io.Writer
	Encoder
	Decoder
}

func (wp *Wrapper) Write(buf []byte) (int, error) {
	if e := wp.Encoder; e != nil {
		buf = e.Encode(buf)
	}
	return wp.w.Write(buf)
}

func (wp *Wrapper) Read(buf []byte) (n int, err error) {
	n, err = wp.r.Read(buf)
	if err != nil {
		return
	}
	if e := wp.Decoder; e != nil {
		n = e.Decode(buf[:n])
	}
	return
}

func (wp *Wrapper) WrapReader(r io.Reader, d Decoder) *Wrapper {
	wp.r = r
	wp.Decoder = d
	return wp
}

func (wp *Wrapper) WrapWriter(w io.Writer, e Encoder) *Wrapper {
	wp.w = w
	wp.Encoder = e
	return wp
}
