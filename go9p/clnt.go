package go9p

import (
	"io"
	"os"
	"go9p.googlecode.com/hg/p"
	"go9p.googlecode.com/hg/p/clnt"
)

const (
	OREAD   = p.OREAD
	OWRITE  = p.OWRITE
	ORDWR   = p.ORDWR
	OEXEC   = p.OEXEC
	OTRUNC  = p.OTRUNC
	OCEXEC  = p.OCEXEC
	ORCLOSE = p.ORCLOSE
)

type Clnt struct {
	*clnt.Clnt
}

// Same as go9p.googlecode.com/hg/p/clnt:Mount, but returning an os.Error
func Mount(net, addr, aname string, user p.User) (*Clnt, os.Error) {
	c, e9 := clnt.Mount(net, addr, aname, user)
	return &Clnt{c}, ToOsError(e9)
}

// Same as go9p.googlecode.com/hg/p/clnt:Clnt.FOpen, but returning a File
// that implements io's Reader Writer and Closer interfaces
func (c *Clnt) FOpen(name string, mode byte) (io.ReadWriteCloser, os.Error) {
	file, e9 := c.Clnt.FOpen(name, mode)
	return &File{file}, ToOsError(e9)
}

type File struct {
	*clnt.File
}

func (f *File) Read(buf []byte) (int, os.Error) {
	n, e9 := f.File.Read(buf)
	return n, ToOsError(e9)
}
func (f *File) Write(buf []byte) (int, os.Error) {
	n, e9 := f.File.Write(buf)
	return n, ToOsError(e9)
}
func (f *File) Close() os.Error {
	return ToOsError(f.File.Close())
}


// convert a 9P error to an os.Error
func ToOsError(e9 *p.Error) os.Error {
	if e9 == nil {
		return nil
	}
	return os.NewError(e9.Error)
}
