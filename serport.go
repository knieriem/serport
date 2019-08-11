package serport

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Port interface {
	SetBaudrate(int) error
	SetParity(byte) error  // odd: 'o', even: 'e', otherwise none
	SetWordlen(int) error  // 5, 6, 7, or 8
	SetStopbits(int) error // 1 or 2

	SetDtr(bool) error
	SetRts(bool) error
	SetRtsCts(bool) error // obey Cts signal, set Rts depending of internal buffer's state

	SetLowLatency(bool) error

	SendBreak(ms int) error

	Delay(ms int)

	// If the Port is remote, after calling Record() the execution of
	// commands will be delayed until Commit() is called.
	Record()
	Commit()

	Drain() error
	Purge(in, out bool)

	Device
}

type Device interface {
	Ctl(cmds ...string) error // accepts commands similar to Plan 9's eia#ctl

	io.ReadWriteCloser
}

type dev struct {
	name     string
	inCtl    bool
	rts, dtr bool
	hw
	encaps Port
}

const (
	initDefault = "b115200 l8 pn r1 s1"
)

func mergeWithDefault(cmds string) string {
	fi := strings.Fields(initDefault)
	f := strings.Fields(cmds)
	r := make([]string, 0, len(fi))
L:
	for _, ci := range fi {
		if ci == "" {
			continue
		}
		for _, c := range f {
			if c[0] == 'D' || c[0] == 'W' {
				break
			}
			if c[0] == ci[0] {
				continue L // exclude c from resulting string
			}
		}
		r = append(r, ci)
	}

	return strings.Join(r, " ") + " " + cmds
}

func (d *dev) Ctl(cmds ...string) error {
	var err error
	updateCtlNow := func() {
		d.inCtl = false
		err = d.updateCtl()
		d.inCtl = true
	}

	d.inCtl = true
	defer func() {
		d.inCtl = false
	}()
	p := d.encaps

	for _, s := range cmds {
		for _, f := range strings.Fields(s) {
			var n int
			var c byte
			var cmd byte

			cmd = f[0]
			if len(f) > 1 {
				if cmd != 'p' {
					n, err = strconv.Atoi(f[1:])
					if err != nil {
						return d.error("ctl", err)
					}
				} else {
					c = f[1]
				}
			}

			//fmt.Printf("Ctl: %c %d\n", cmd, n)
			switch cmd {
			case 'd':
				err = p.SetDtr(n == 1)
			case 'r':
				err = p.SetRts(n == 1)
			case 'm':
				err = p.SetRtsCts(n != 0)
			case 'D':
				updateCtlNow()
				p.Delay(n)
			case 'W':
				updateCtlNow()
				if err != nil {
					break
				}
				_, err = p.Write([]byte{byte(n)})
			case 'b':
				err = p.SetBaudrate(n)
			case 'l':
				err = p.SetWordlen(n)
			case 'k':
				err = p.SendBreak(n)
			case 'p':
				err = p.SetParity(c)
			case 's':
				err = p.SetStopbits(n)
			case 'L':
				err = p.SetLowLatency(n != 0)
			default:
				err = d.errorf("ctl", "unknown command: %q", string(cmd))
			}
			if err != nil {
				return err
			}
		}
	}
	d.inCtl = false
	return d.updateCtl()
}

func (d *dev) Delay(ms int) {
	d.Drain()
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (d *dev) Record() {
}

func (d *dev) Commit() {
}

func (d *dev) error(action string, err error) error {
	return pathError(action, d.name, err)
}

func (d *dev) errorf(action string, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	return pathError(action, d.name, err)
}

func pathError(op, path string, err error) error {
	return &os.PathError{Op: op, Path: path, Err: err}
}

type LineState struct {
	Csr  bool
	Dsr  bool
	Ring bool
	Dcd  bool
}

func SetEncapsulatingPort(dest, p Port) {
	if dev, ok := dest.(*dev); ok {
		dev.encaps = p
	}
}
