// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived from pnm/reader.go, which is based on
// the structure of image/gif/reader.go.

// Package xbm implements an X11 Bitmap decoder.
package xbm

import (
	"bufio"
	"errors"

	img "github.com/knieriem/g/image"
	"image"
	"io"
	"strconv"
	"strings"
)

// If the io.Reader does not also have ReadLine, then decode will introduce its own buffering.
type reader interface {
	io.Reader
	ReadLine() (line []byte, isPrefix bool, err error)
}

// decoder is the type used to decode an XBM file.
type decoder struct {
	r reader

	// From header.
	width  int
	height int
	line   []byte
}

// decode reads an X11 bitmap from r and stores the result in d.
func (d *decoder) decode(r io.Reader, configOnly bool) (im *img.Bitmap, err error) {
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok {
				err = errors.New("xbm:" + s)
			} else {
				err = r.(error)
			}
		}
	}()

	// Add buffering if r does not provide ReadByte.
	if rr, ok := r.(reader); ok {
		d.r = rr
	} else {
		d.r = bufio.NewReader(r)
	}

	d.readHeader()
	if configOnly {
		return
	}

	im = img.NewBitmap(d.width, d.height)
	p := im.Pix

	if len(d.line) > 0 {
		p = d.scan(p, d.line)
	}
	for len(p) > 0 {
		line, isPfx, err := d.r.ReadLine()
		if isPfx {
			goto malformed
		}
		if err != nil {
			return nil, err
		}
		p = d.scan(p, line)
	}
	return

malformed:
	return nil, errors.New("xbm: data probably malformed")
}

func (d *decoder) readHeader() {
	var name string

	for {
		line, isPfx, err := d.r.ReadLine()
		if isPfx {
			goto malformed
		}
		if err != nil {
			panic(err)
		}
		linestr := string(line)
		f := strings.Fields(linestr)
		if len(f) < 3 {
			continue
		}
		switch f[0] {
		case "#define":
			if iu := strings.LastIndex(f[1], "_"); iu != -1 {
				s := f[1][:iu]
				switch name {
				case "":
					name = s
				default:
					panic("name mismatch")
				case s:
				}
				val, _ := strconv.Atoi(f[2])
				switch f[1][iu+1:] {
				case "width":
					d.width = val
				case "height":
					d.height = val
				}
			}
		case "static":
			if d.width == 0 || d.height == 0 {
				goto malformed
			}
			if f[1] != "char" && f[2] != "char" {
				panic("data type not supported")
			}
			if bi := strings.Index(linestr, "{"); bi == -1 {
				goto malformed
			} else {
				d.line = line[bi+1:]
			}
			return
		}
	}
malformed:
	panic("probably malformed header")
}

func (d *decoder) scan(dst []byte, line []byte) []byte {
	var ndst = len(dst)
	var n = len(line)
	var di = 0

	for i := 0; i < n && di != ndst; i++ {
		if line[i] == 'x' {
			if n-i < 3 {
				panic("malformed data")
			}
			dst[di] = unhex(line[i+2])<<4 | unhex(line[i+1])
			di++
		}
	}
	return dst[di:]
}

var flipped = []byte{0, 8, 4, 0xC, 2, 0xA, 6, 0xE, 1, 9, 5, 0xD, 3, 0xB, 7, 0xF}

func unhex(h byte) (b uint8) {
	switch {
	case h >= '0' && h <= '9':
		b = h - '0'
	case h >= 'A' && h <= 'F':
		b = 10 + h - 'A'
	case h >= 'a' && h <= 'f':
		b = 10 + h - 'a'
	default:
		panic("malformed data")
	}
	return flipped[b]
}

// Decode reads an XBM image from r and returns the first embedded
// image as an image.Image.
func Decode(r io.Reader) (im image.Image, err error) {
	var d decoder
	return d.decode(r, false)
}

// DecodeConfig returns the color model and dimensions of an XBM image
// without decoding the entire image.
func DecodeConfig(r io.Reader) (ic image.Config, err error) {
	var d decoder
	if _, err = d.decode(r, true); err == nil {
		ic = image.Config{img.BinaryColorModel, d.width, d.height}
	}
	return
}

func init() {
	image.RegisterFormat("xbm", "/*", Decode, DecodeConfig)
	image.RegisterFormat("xbm", "#defi", Decode, DecodeConfig)
}
