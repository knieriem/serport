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
	Ctl(cmds ...string) error // accepts commands similar to Plan 9's eia#ctl

	SetBaudrate(int) error
	SetParity(byte) error  // odd: 'o', even: 'e', otherwise none
	SetWordlen(int) error  // 5, 6, 7, or 8
	SetStopbits(int) error // 1 or 2

	SetDtr(bool) error
	SetRts(bool) error
	SetRtsCts(bool) error // obey Cts signal, set Rts depending of internal buffer's state

	SendBreak(ms int) error

	Delay(ms int)

	// If the Port is remote, after calling Record() the execution of
	// commands will be delayed until Commit() is called.
	Record()
	Commit()

	Drain() error
	Purge(in, out bool)

	io.ReadWriteCloser
}

type dev struct {
	name     string
	inCtl    bool
	rts, dtr bool
	hw
	encaps Port
}

func (p *dev) Ctl(cmds ...string) error {
	var err error
	updateCtlNow := func() {
		p.inCtl = false
		err = p.updateCtl()
		p.inCtl = true
	}

	p.inCtl = true
	defer func() {
		p.inCtl = false
	}()
	d := p.encaps

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
						return p.error("ctl", err)
					}
				} else {
					c = f[1]
				}
			}

			//fmt.Printf("Ctl: %c %d\n", cmd, n)
			switch cmd {
			case 'd':
				err = d.SetDtr(n == 1)
			case 'r':
				err = d.SetRts(n == 1)
			case 'm':
				err = d.SetRtsCts(n != 0)
			case 'D':
				updateCtlNow()
				d.Delay(n)
			case 'W':
				updateCtlNow()
				if err != nil {
					break
				}
				_, err = d.Write([]byte{byte(n)})
			case 'b':
				err = d.SetBaudrate(n)
			case 'l':
				err = d.SetWordlen(n)
			case 'k':
				err = d.SendBreak(n)
			case 'p':
				err = d.SetParity(c)
			case 's':
				err = d.SetStopbits(n)
			}
			if err != nil {
				return err
			}
		}
	}
	p.inCtl = false
	return p.updateCtl()
}

func (d *dev) Delay(ms int) {
	d.Drain()
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (d *dev) Record() {
}

func (d *dev) Commit() {
}

func (p *dev) error(action string, err error) error {
	return pathError(action, p.name, err)
}

func (p *dev) errorf(action string, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	return pathError(action, p.name, err)
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
