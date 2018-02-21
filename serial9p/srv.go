/*
Package serial9p makes serial ports accessible through the 9P protocol.
A serport.Port may be exported as a 9P file server that should
behave much like Plan 9's /dev/eia* serial line devices.
Functions MountConn and OpenFSDev provide 9P client access.

See http://plan9.bell-labs.com/magic/man2html/3/uart for details about
the file system interface, and cmd/sercom for an example.
*/
package serial9p

import (
	"strings"
	"sync"

	"github.com/lionkov/go9p/p"
	"github.com/lionkov/go9p/p/srv"

	"github.com/knieriem/g/ioutil"
	"github.com/knieriem/serport"
)

type ctl struct {
	file
	dev           serport.Port
	dataUnblockch chan bool
	record        bool
	rlist         []string
}
type data struct {
	file
	m         sync.Mutex
	dev       serport.Port
	rch       ioutil.RdChannels
	fid       *srv.Fid
	unblockch chan bool
}

type file struct {
	srv.File
}

func (*file) Wstat(*srv.FFid, *p.Dir) error {
	return nil
}

func (c *ctl) Write(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	var err error

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
	return len(buf), err
}

func (d *data) Read(fid *srv.FFid, buf []byte, offset uint64) (n int, err error) {
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
	return
}
func (d *data) Write(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	return d.dev.Write(buf)
}

func (d *data) Clunk(f *srv.FFid) error {
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

// Link ctl and data files (that wrap a previously opened
// serial port) into an existing 9P file tree `dir'.
func RegisterFiles9P(dir *srv.File, dev serport.Port, user p.User) (err error) {
	c := new(ctl)
	c.dev = dev
	err = c.Add(dir, "ctl", user, nil, 0666, c)
	if err != nil {
		return
	}

	d := new(data)
	d.dev = dev
	d.rch = ioutil.ChannelizeReader(dev, nil)
	d.unblockch = make(chan bool)
	c.dataUnblockch = d.unblockch
	err = d.Add(dir, "data", user, nil, 0666, d)

	return
}
