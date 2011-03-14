package sercom

import (
	"io"
	"os"
	"strconv"
	p "github.com/knieriem/g/go9p"
)

type dev9 struct {
	clnt *p.Clnt
	data io.ReadWriteCloser
	ctl	io.WriteCloser
	recording bool
}

// Connect to a 9P server that is listening at `addr',
// and wrap ctl and data files of a remote serial device
// into a Port.
// Basename is the a file name prefix or directory,
// where the device is expected to be found. It can
// be "", if ctl and data files live in the 9P servers
// root directory
func Connect9P(addr, basename string) (port Port, err os.Error) {
	c, err := p.Mount("tcp", addr, "", p.CurrentUser())
	if err != nil {
		return
	}
	d := new(dev9)
	d.clnt = c
	if d.data, err = d.clnt.FOpen(basename+"/data", p.ORDWR); err == nil {
		if d.ctl, err = d.clnt.FOpen(basename+"/ctl", p.OWRITE); err != nil {
			goto noctl
		}
	} else if d.data, err = d.clnt.FOpen(basename, p.ORDWR); err == nil {
		if d.ctl, err = d.clnt.FOpen(basename+"ctl", p.OWRITE); err != nil {
			goto noctl
		}
	} else {
		goto unmount
	}
	port = d
	return

noctl:
	d.data.Close()
unmount:
	d.clnt.Unmount()
	return
}

func (d *dev9) Read(buf []byte) (n int, err os.Error) {
	return d.data.Read(buf)
}

func (d *dev9) Write(buf []byte) (n int, err os.Error) {
	if d.recording {
		for _, b := range buf {
			if err = d.cmdi('W', int(b)); err != nil {
				break
			}
		}
		return
	}
	return d.data.Write(buf)
}

func (d *dev9) Close() os.Error {
	d.ctl.Close()
	d.data.Close()
	d.clnt.Unmount()
	return nil
}

func (d *dev9) Drain() os.Error {
	return nil
}

func (d *dev9) Purge(in, out bool) {
}

func (d *dev9) Delay(ms int) {
	d.cmdi('D', ms)
}

func (d *dev9) Record() {
	d.cmd("{")
	d.recording = true
}

func (d *dev9) Commit() {
	d.cmd("}")
	d.recording = false
}

func (d *dev9) Ctl(cmd string) os.Error {
	return d.cmd(cmd)
}

func (d *dev9) SetBaudrate(val int) (err os.Error) {
	return d.cmdi('b', val)
}

func (d *dev9) SetWordlen(n int) os.Error {
	return d.cmdi('l', n)
}

func (d *dev9) SetParity(parity byte) os.Error {
	return d.cmd("p" + string(parity))
}

func (d *dev9) SetStopbits(n int) (err os.Error) {
	if n==1 || n==2 {
		return d.cmdi('s', n)
	}
	return os.NewError("invalid number of stopbits")
}

func (d *dev9) SetRts(on bool) os.Error {
	return d.cmdbool('r', on)
}
func (d *dev9) SetDtr(on bool) os.Error {
	return d.cmdbool('d', on)
}

func (d *dev9) SetRtsCts(on bool) os.Error {
	return d.cmdbool('m', on)
}

func (d *dev9) cmd(c string) (err os.Error) {
	_, err = d.ctl.Write([]byte(c))
	return err
}

func (d *dev9) cmdbool(c byte, on bool) (err os.Error) {
	var msg = []byte{c, '0', c, '1'}

	if on {
		msg = msg[2:]
	}
	_, err = d.ctl.Write(msg[:2])
	return
}

func (d *dev9) cmdi(c byte, val int) (err os.Error) {
	return d.cmd(string(c)+strconv.Itoa(val))
}
