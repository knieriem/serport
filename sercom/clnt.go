package sercom

import (
	"io"
	"os"
	"strconv"
	p "github.com/knieriem/g/go9p"
)

type dev9 struct {
	clnt *p.Clnt
	fdev
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

func (d *dev9) Close() os.Error {
	d.fdev.Close()
	d.clnt.Unmount()
	return nil
}


type fdev struct {
	data      io.ReadWriteCloser
	ctl       io.WriteCloser
	recording bool
}

// Connect to a serial device that is accessible in
// the local file system, e.g. driven by a 9pfuse mounted
// 9P service.
// Devdir is the name of a directory, where the files
// "ctl" and "data" are expected to be found.
func OpenFsDev(devdir string) (port Port, err os.Error) {
	d := new(fdev)
	if d.data, err = os.OpenFile(devdir+"/data", os.O_RDWR, 0); err == nil {
		if d.ctl, err = os.OpenFile(devdir+"/ctl", os.O_RDWR, 0); err != nil {
			d.data.Close()
		} else {
			port = d
		}
	}
	return
}

func (d *fdev) Read(buf []byte) (n int, err os.Error) {
	return d.data.Read(buf)
}

func (d *fdev) Write(buf []byte) (n int, err os.Error) {
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

func (d *fdev) Close() os.Error {
	d.cmd("U") // unlock remote Read()
	d.ctl.Close()
	d.data.Close()
	return nil
}

func (d *fdev) Drain() os.Error {
	return nil
}

func (d *fdev) Purge(in, out bool) {
}

func (d *fdev) Delay(ms int) {
	d.cmdi('D', ms)
}

func (d *fdev) Record() {
	d.cmd("{")
	d.recording = true
}

func (d *fdev) Commit() {
	d.cmd("}")
	d.recording = false
}

func (d *fdev) Ctl(cmds ...string) (err os.Error) {
	for _, s := range cmds {
		if err = d.cmd(s); err != nil {
			break
		}
	}
	return
}

func (d *fdev) SetBaudrate(val int) (err os.Error) {
	return d.cmdi('b', val)
}

func (d *fdev) SetWordlen(n int) os.Error {
	return d.cmdi('l', n)
}

func (d *fdev) SetParity(parity byte) os.Error {
	return d.cmd("p" + string(parity))
}

func (d *fdev) SetStopbits(n int) (err os.Error) {
	if n == 1 || n == 2 {
		return d.cmdi('s', n)
	}
	return os.NewError("invalid number of stopbits")
}

func (d *fdev) SetRts(on bool) os.Error {
	return d.cmdbool('r', on)
}
func (d *fdev) SetDtr(on bool) os.Error {
	return d.cmdbool('d', on)
}

func (d *fdev) SetRtsCts(on bool) os.Error {
	return d.cmdbool('m', on)
}

func (d *fdev) cmd(c string) (err os.Error) {
	_, err = d.ctl.Write([]byte(c))
	return err
}

func (d *fdev) cmdbool(c byte, on bool) (err os.Error) {
	var msg = []byte{c, '0', c, '1'}

	if on {
		msg = msg[2:]
	}
	_, err = d.ctl.Write(msg[:2])
	return
}

func (d *fdev) cmdi(c byte, val int) (err os.Error) {
	return d.cmd(string(c) + strconv.Itoa(val))
}
