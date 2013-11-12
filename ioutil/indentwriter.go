package ioutil

import "io"

// An IndentWriter indents each '\n'-terminated line by
// a number of '\t' characters.
// Blank lines can be indented too, optionally.
//
type IndentWriter struct {
	io.Writer
	indent []byte
	inLine bool

	indentBlank bool
}

func NewIndentWriter(w io.Writer, pfx []byte) (iw *IndentWriter) {
	iw = new(IndentWriter)
	iw.Writer = w
	iw.indent = pfx
	return
}

func (w *IndentWriter) Write(buf []byte) (n int, err error) {
	var i0 = 0

	for i := range buf {
		if buf[i] == '\n' {
			if i == i0 || buf[i-1] == '\r' {
				if !w.indentBlank {
					goto skipIndent
				}
			}
			if w.inLine {
				w.inLine = false
			} else if _, err = w.Writer.Write(w.indent); err != nil {
				return
			}

		skipIndent:
			if _, err = w.Writer.Write(buf[i0 : i+1]); err != nil {
				return
			}
			i0 = i + 1
		}
	}

	n = len(buf)
	if i0 != n {
		w.inLine = true
		_, err = w.Writer.Write(buf[i0:])
	}
	return
}

func (w *IndentWriter) IndentBlankLines(val bool) *IndentWriter {
	w.indentBlank = val
	return w
}

func (w *IndentWriter) SubWriter() *IndentWriter {
	sub := *w
	return (&sub).Inc()
}

func (w *IndentWriter) Inc() *IndentWriter {
	w.indent = append(w.indent, '\t')
	return w
}

func (w *IndentWriter) Dec() *IndentWriter {
	if n := len(w.indent); n > 0 {
		w.indent = w.indent[:n-1]
	}
	return w
}
