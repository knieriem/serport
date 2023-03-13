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

	SendBreak(time.Duration) error

	Delay(time.Duration)

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
	// StdConf specifies a port configuration that may be used
	// as inictl argument to [Open].
	StdConf = "b115200 l8 pn r1 s1"
)

// MergeCtlCmds concatenates control cmd strings, skipping commands
// from preceding strings if they occur, possibly with different arguments,
// in subsequent strings too.
func MergeCtlCmds(ctlCmds ...string) string {
	switch len(ctlCmds) {
	case 0:
		return ""
	case 1:
		return mergeCtlCmds(StdConf, ctlCmds[0])
	}
	cmds := ctlCmds[0]
	for _, cmds2 := range ctlCmds[1:] {
		cmds = mergeCtlCmds(cmds, cmds2)
	}
	return cmds
}

func mergeCtlCmds(cmds1, cmds2 string) string {
	fi := strings.Fields(cmds1)
	f := strings.Fields(cmds2)
	r := make([]string, 0, len(fi))

	for _, ci := range fi {
		if ci == "" {
			continue
		}
		for i, c := range f {
			if c == "" {
				continue
			}
			if c[0] == 'D' || c[0] == 'W' {
				break
			}
			if c[0] == ci[0] {
				ci = c
				f[i] = ""
			}
		}
		r = append(r, ci)
	}

	fOrig, f := f, f[:0]
	for _, c := range fOrig {
		if c == "" {
			continue
		}
		f = append(f, c)
	}
	r = append(r, f...)
	return strings.Join(r, " ")
}

var ctlNamespaceMap = map[string]*ctlNamespace{
	"std": stdNamespace,
}

type ctlNamespace struct {
	runCmd    ctlRunFunc
	updateDrv func(*dev) error
	charCmds  string
}

type ctlRunFunc func(d *dev, cmd, c byte, n int) error

func (d *dev) Ctl(cmds ...string) error {
	var err error

	defer func() {
		d.inCtl = false
	}()

	subNsID := ""
	nsStd := ctlNamespaceMap["std"]
	ns := nsStd
	for _, s := range cmds {
		for _, f := range strings.Fields(s) {
			var n int
			var c byte

			d.inCtl = true
			cmd := f[0]
			if len(f) > 1 {
				arg := f[1:]
				if cmd == '.' {
					if arg[0] == '.' {
						if subNsID == "" {
							return d.errorf("ctl", "found \"..\", but sub-namespace not active")
						}
						arg = arg[1:]
					} else if iDot := strings.IndexByte(arg, '.'); iDot == -1 {
						return d.errorf("ctl", "missing dot after namespace id")
					} else {
						id := arg[:iDot]
						arg = arg[iDot+1:]
						if subNsID != id {
							if subNsID != "" {
								err = ns.updateDrv(d)
								if err != nil {
									return err
								}
							}
							nsNew, ok := ctlNamespaceMap[id]
							if !ok {
								return d.errorf("ctl", "namespace not implemented: %q", id)
							}
							ns = nsNew
							subNsID = id
						}
					}
					if len(arg) == 0 {
						return d.errorf("ctl", "command missing, found: %q", f)
					}
					cmd = arg[0]
					arg = arg[1:]
				} else if subNsID != "" {
					err = ns.updateDrv(d)
					if err != nil {
						return err
					}
					subNsID = ""
					ns = nsStd
				}
				if strings.IndexByte(ns.charCmds, cmd) == -1 {
					n, err = strconv.Atoi(arg)
					if err != nil {
						return d.error("ctl", err)
					}
				} else {
					c = arg[0]
				}
			}
			err := ns.runCmd(d, cmd, c, n)
			if err != nil {
				return err
			}
		}
	}
	if subNsID != "" {
		err = ns.updateDrv(d)
		if err != nil {
			return err
		}
	}

	if !d.inCtl {
		return nil
	}
	d.inCtl = false
	return d.updateCtl()
}

func (d *dev) updateCtlNow() error {
	saved := d.inCtl
	d.inCtl = false
	defer func() { d.inCtl = saved }()
	return d.updateCtl()
}

var stdNamespace = &ctlNamespace{
	runCmd: func(d *dev, cmd, c byte, n int) error {
		p := d.encaps
		switch cmd {
		case 'd':
			return p.SetDtr(n == 1)
		case 'r':
			return p.SetRts(n == 1)
		case 'm':
			return p.SetRtsCts(n != 0)
		case 'D':
			err := d.updateCtlNow()
			if err != nil {
				return err
			}
			p.Delay(time.Duration(n) * time.Millisecond)
			return nil
		case 'W':
			err := d.updateCtlNow()
			if err != nil {
				return err
			}
			_, err = p.Write([]byte{byte(n)})
			return err
		case 'b':
			return p.SetBaudrate(n)
		case 'l':
			return p.SetWordlen(n)
		case 'k':
			return p.SendBreak(time.Duration(n) * time.Millisecond)
		case 'p':
			return p.SetParity(c)
		case 's':
			return p.SetStopbits(n)
		case 'L':
			return p.SetLowLatency(n != 0)
		default:
			return d.errorf("ctl", "unknown command: %q", string(cmd))
		}
	},
	charCmds: "p",
}

func (d *dev) Delay(duration time.Duration) {
	d.Drain()
	time.Sleep(duration)
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
