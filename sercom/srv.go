package sercom

import (
	"fmt"
	"os"
	"sync"

	"go9p.googlecode.com/hg/p"
	"go9p.googlecode.com/hg/p/srv"
	"github.com/knieriem/g/go9p"
)

var (
	Debug    bool
	Debugall bool
)

type ctl struct {
	srv.File
	dev Port
}
type data struct {
	m	sync.Mutex
	srv.File
	dev     Port
	clunked bool
	tmp     []byte
}

func (c *ctl) Write(fid *srv.FFid, buf []byte, offset uint64) (int, *p.Error) {
	err := c.dev.Ctl(string(buf))
	return len(buf), go9p.ToError(err)
}

func (d *data) Read(fid *srv.FFid, buf []byte, offset uint64) (n int, e9 *p.Error) {
	var err os.Error

	d.m.Lock()
	defer d.m.Unlock()

	if nt := len(d.tmp); nt != 0 {
		n = len(buf)
		if n > nt {
			n = nt
		}
		copy(buf, d.tmp[:n])
		d.tmp = d.tmp[n:]
	} else {
		n, err = d.dev.Read(buf)
		if d.clunked {
			d.tmp = make([]byte, n)
			copy(d.tmp, buf[:n])
			d.clunked = false
			n = 0
		}
	}
	e9 = go9p.ToError(err)
	return
}
func (d *data) Write(fid *srv.FFid, buf []byte, offset uint64) (int, *p.Error) {
	n, err := d.dev.Write(buf)
	return n, go9p.ToError(err)
}
func (d *data) Clunk(*srv.FFid) *p.Error {
	d.clunked = true
	return nil
}

// Serve a previously opened serial device via 9P.
// `addr' shoud be of form "host:port", where host
// may be missing.
func Serve9P(addr string, dev Port) os.Error {
	user := go9p.CurrentUser()
	root := new(srv.File)
	err := root.Add(nil, "/", user, nil, p.DMDIR|0555, nil)
	if err != nil {
		goto error
	}

	c := new(ctl)
	c.dev = dev
	err = c.Add(root, "ctl", user, nil, 0664, c)
	if err != nil {
		goto error
	}
	d := new(data)
	d.dev = dev
	err = d.Add(root, "data", user, nil, 0664, d)
	if err != nil {
		goto error
	}

	s := srv.NewFileSrv(root)
	s.Dotu = true

	switch {
	case Debugall:
		s.Debuglevel = 2
	case Debug:
		s.Debuglevel = 1
	}

	s.Start(s)
	err = s.StartNetListener("tcp", addr)
	if err != nil {
		goto error
	}

	return nil

error:
	return os.NewError(fmt.Sprintf("Error: %s %d", err.Error, err.Errornum))
}
