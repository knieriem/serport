package sercom

import (
	"os"
	"strings"
	"sync"

	"go9p.googlecode.com/hg/p"
	"go9p.googlecode.com/hg/p/srv"
	"github.com/knieriem/g/go9p"
	"github.com/knieriem/g/ioutil"
)

var (
	Debug    bool
	Debugall bool
)

type ctl struct {
	file
	dev           Port
	dataUnblockch chan bool
	record        bool
	rlist         []string
}
type data struct {
	file
	m         sync.Mutex
	dev       Port
	rch       ioutil.RdChannels
	fid       *srv.Fid
	unblockch chan bool
}

type file struct {
	srv.File
}

func (*file) Wstat(*srv.FFid, *p.Dir) *p.Error {
	return nil
}

func (c *ctl) Write(fid *srv.FFid, buf []byte, offset uint64) (int, *p.Error) {
	var err os.Error

	for _, cmd := range strings.Fields(string(buf)) {
		switch cmd {
		case "U":
			select {
			case c.dataUnblockch <- true:
			default:
			}
		case "{":
			c.record = true
		case "}":
			c.record = false
			if len(c.rlist) != 0 {
				err = c.dev.Ctl(c.rlist...)
				c.rlist = c.rlist[:0]
			}
		default:
			if c.record {
				c.rlist = append(c.rlist, cmd)
			} else {
				err = c.dev.Ctl(cmd)
			}
		}
		if err != nil {
			break
		}
	}
	return len(buf), go9p.ToError(err)
}

func (d *data) Read(fid *srv.FFid, buf []byte, offset uint64) (n int, e9 *p.Error) {
	var err os.Error

	d.m.Lock()
	defer d.m.Unlock()

	d.fid = fid.Fid
	select {
	case in := <-d.rch.Data:
		n = len(in.Data)
		if n > len(buf) {
			n = len(buf)
		}
		copy(buf, in.Data[:n])
		d.rch.Req <- n
		err = in.Err
	case <-d.unblockch:
		n = 0
	}
	d.fid = nil
	e9 = go9p.ToError(err)
	return
}
func (d *data) Write(fid *srv.FFid, buf []byte, offset uint64) (int, *p.Error) {
	n, err := d.dev.Write(buf)
	return n, go9p.ToError(err)
}

func (d *data) Clunk(f *srv.FFid) *p.Error {
	if f.Fid == d.fid {
		// The fid is clunked while a Tread is outstanding.
		// Unblock data.Read() so that the Rread can be sent.
		select {
		case d.unblockch <- true:
		default:
		}
	}
	return nil
}

// Serve a previously opened serial device via 9P.
// `addr' shoud be of form "host:port", where host
// may be missing.
func Serve9P(addr string, dev Port) (err os.Error) {
	user := go9p.CurrentUser()
	root := new(srv.File)
	err = root.Add(nil, "/", user, nil, p.DMDIR|0555, nil)
	if err != nil {
		return
	}

	c := new(ctl)
	c.dev = dev
	err = c.Add(root, "ctl", user, nil, 0664, c)
	if err != nil {
		return
	}

	d := new(data)
	d.dev = dev
	d.rch = ioutil.ChannelizeReader(dev, nil)
	d.unblockch = make(chan bool)
	c.dataUnblockch = d.unblockch
	err = d.Add(root, "data", user, nil, 0664, d)
	if err != nil {
		return
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
	return
}
