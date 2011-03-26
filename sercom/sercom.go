package sercom

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Port interface {
	Ctl(cmds ...string) os.Error // accepts commands similar to Plan 9's eia#ctl

	SetBaudrate(int) os.Error
	SetParity(byte) os.Error  // odd: 'o', even: 'e', otherwise none
	SetWordlen(int) os.Error  // 5, 6, 7, or 8
	SetStopbits(int) os.Error // 1 or 2

	SetDtr(bool) os.Error
	SetRts(bool) os.Error
	SetRtsCts(bool) os.Error // obey Cts signal, set Rts depending of internal buffer's state

	Delay(ms int)

	// If the Port is remote, after calling Record() the execution of
	// commands will be delayed until Commit() is called.
	Record()
	Commit()

	Drain() os.Error
	Purge(in, out bool)

	io.ReadWriteCloser
}


type dev struct {
	name     string
	inCtl    bool
	rts, dtr bool
	hw
}

func (d *dev) Ctl(cmds ...string) os.Error {
	updateCtlNow := func() {
		d.inCtl = false
		d.updateCtl()
		d.inCtl = true
	}

	d.inCtl = true

	for _, s := range cmds {
		for _, f := range strings.Fields(s) {
			var n int
			var c byte
			var cmd byte

			switch len(f) {
			default:
				n, _ = strconv.Atoi(f[1:])
				c = f[1]
				fallthrough
			case 1:
				cmd = f[0]
			}
			//fmt.Printf("Ctl: %c %d\n", cmd, n)
			switch cmd {
			case 'd':
				d.SetDtr(n == 1)
			case 'r':
				d.SetRts(n == 1)
			case 'm':
				d.SetRtsCts(n != 0)
			case 'D':
				updateCtlNow()
				d.Delay(n)
			case 'W':
				updateCtlNow()
				d.Write([]byte{byte(n)})
			case 'b':
				d.SetBaudrate(n)
			case 'l':
				d.SetWordlen(n)
			case 'p':
				d.SetParity(c)
			case 's':
				d.SetStopbits(n)
			}
		}
	}
	d.inCtl = false
	return d.updateCtl()
}

func (d *dev) Delay(ms int) {
	d.Drain()
	time.Sleep(int64(ms) * 1e6)
}

func (d *dev) Record() {
}

func (d *dev) Commit() {
}

func (p *dev) errno(action string, e int) os.Error {
	return &os.PathError{action, p.name, os.Errno(e)}
}

func (p *dev) errorf(action string, format string, args ...interface{}) os.Error {
	err := os.NewError(fmt.Sprintf(format, args...))
	return &os.PathError{action, p.name, err}
}

type LineState struct {
	Csr  bool
	Dsr  bool
	Ring bool
	Dcd  bool
}
