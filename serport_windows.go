package serport

import (
	win "github.com/knieriem/g/syscall"
	"io"
	"path/filepath"
	"sync"
	"syscall"
)

const (
	initDefault = "b115200 l8 pn r1 s1"
)

type hw struct {
	fd syscall.Handle
	sync.Mutex
	closing bool

	initDone    bool
	dcb, dcbsav win.DCB

	ev struct {
		w syscall.Handle
		r syscall.Handle
	}
}

func Open(file string, inictl string) (p Port, err error) {
	const (
		access     = syscall.GENERIC_READ | syscall.GENERIC_WRITE
		sharemode  = 0
		createmode = syscall.OPEN_EXISTING
		flags      = win.FILE_FLAG_OVERLAPPED
	)

	// make sure COM interfaces with numbers >9 get prefixed properly
	if match, _ := filepath.Match("[cC][oO][mM]1[0-9]", file); match {
		file = `\\.\` + file
	}

	fd, e := syscall.CreateFile(syscall.StringToUTF16Ptr(file), access, sharemode, nil, createmode, flags, 0)
	if e != nil {
		goto error
	}

	goto try
error:
	err = pathError("open", file, e)
	return

try:
	d := new(dev)
	d.fd = fd
	d.name = file
	d.encaps = d

	if err = d.Ctl(initDefault, inictl); err != nil {
		return
	}
	d.initDone = true

	if d.ev.r, e = win.CreateEvent(win.EvManualReset, !win.EvInitiallyOn); e != nil {
		goto error
	}
	if d.ev.w, e = win.CreateEvent(win.EvManualReset, !win.EvInitiallyOn); e != nil {
		goto error
	}

	cto := win.CommTimeouts{
		ReadIntervalTimeout:        ^uint32(0),
		ReadTotalTimeoutMultiplier: ^uint32(0),
		ReadTotalTimeoutConstant:   200,
		//	ReadIntervalTimeout: 10,
	}
	if e = win.SetCommTimeouts(d.fd, &cto); e != nil {
		goto error
	}
	if e = win.SetupComm(d.fd, 4096, 4096); e != nil {
		goto error
	}
	p = d
	return
}

func (p *dev) Read(buf []byte) (int, error) {
	var done uint32

	p.Lock()
	defer p.Unlock()
	for {
		var ov syscall.Overlapped

		if p.closing {
			return 0, io.EOF
		}
		ov.HEvent = p.ev.r
		if e := syscall.ReadFile(p.fd, buf, &done, &ov); e != nil {
			if e != syscall.ERROR_IO_PENDING {
				return 0, p.error("reading from", e)
			}
			if e = win.GetOverlappedResult(p.fd, &ov, &done, 1); e != nil {
				return 0, p.error("reading from", e)
			}
		}
		if done > 0 {
			break
		}
	}
	return int(done), nil
}

func (p *dev) Write(buf []byte) (int, error) {
	var done uint32

	for {
		var ov syscall.Overlapped

		ov.HEvent = p.ev.w
		if e := syscall.WriteFile(p.fd, buf, &done, &ov); e != nil {
			if e != syscall.ERROR_IO_PENDING {
				return 0, p.error("writing to", e)
			}
			if e = win.GetOverlappedResult(p.fd, &ov, &done, 1); e != nil {
				return 0, p.error("writing to", e)
			}
		}
		if done > 0 {
			break
		}
	}
	return int(done), nil
}

func (d *dev) Close() (err error) {
	d.closing = true
	d.Lock()
	defer d.Unlock()
	d.Drain()
	syscall.CloseHandle(d.ev.r)
	syscall.CloseHandle(d.ev.w)
	if e := syscall.CloseHandle(d.fd); e != nil {
		err = d.error("close", e)
	}
	return nil
}

func (d *dev) Drain() (err error) {
	if e := win.FlushFileBuffers(d.fd); e != nil {
		err = d.error("drain", e)
	}
	return
}

func (d *dev) Purge(in, out bool) {
	// TBD
}

func (d *dev) SetBaudrate(val int) error {
	d.dcb.BaudRate = uint32(val)
	return d.updateCtl()
}

func (d *dev) SetWordlen(n int) error {
	switch n {
	case 5, 6, 7, 8:
		d.dcb.ByteSize = uint8(n)
	}
	return d.updateCtl()
}

func (d *dev) SetParity(val byte) error {
	p := &d.dcb.Parity
	switch val {
	case 'o':
		*p = win.ODDPARITY
	case 'e':
		*p = win.EVENPARITY
	default:
		*p = win.NOPARITY
	}
	return d.updateCtl()
}

func (d *dev) SetStopbits(n int) error {
	switch n {
	case 1:
		d.dcb.StopBits = win.ONESTOPBIT
	case 2:
		d.dcb.StopBits = win.TWOSTOPBITS
	default:
		return d.errorf("open", "invalid number of stopbits: %d", n)
	}
	return d.updateCtl()
}

func (d *dev) SetRts(on bool) (err error) {
	d.rts = on
	setRtsFlags(&d.dcb, on)
	if !d.initDone {
		return
	}
	setRtsFlags(&d.dcbsav, on) // fake
	if on {
		return d.commfn("set rts", win.SETRTS)
	}
	return d.commfn("clr rts", win.CLRRTS)
}

func (d *dev) SetDtr(on bool) (err error) {
	d.dtr = on
	setDtrFlags(&d.dcb, on)
	if !d.initDone {
		return
	}
	setDtrFlags(&d.dcbsav, on) // fake
	if on {
		return d.commfn("set dtr", win.SETDTR)
	}
	return d.commfn("clr dtr", win.CLRDTR)
}

func (d *dev) commfn(name string, f int) (err error) {
	if e := win.EscapeCommFunction(d.fd, uint32(f)); e != nil {
		return d.error(name, e)
	}
	return
}

func (d *dev) SetRtsCts(on bool) error {
	dcb := &d.dcb

	if on {
		dcb.Flags &^= win.DCBmRtsControl << win.DCBpRtsControl
		dcb.Flags |= win.DCBfOutxCtsFlow
		dcb.Flags |= win.RTS_CONTROL_HANDSHAKE << win.DCBpRtsControl
	} else {
		dcb.Flags &^= win.DCBfOutxCtsFlow
		setRtsFlags(dcb, d.rts)
	}
	return d.updateCtl()
}

func setRtsFlags(dcb *win.DCB, on bool) {
	dcb.Flags &^= win.DCBmRtsControl << win.DCBpRtsControl
	if on {
		dcb.Flags |= win.RTS_CONTROL_ENABLE << win.DCBpRtsControl
	} else {
		dcb.Flags |= win.RTS_CONTROL_DISABLE << win.DCBpRtsControl
	}
}
func setDtrFlags(dcb *win.DCB, on bool) {
	dcb.Flags &^= win.DCBmDtrControl << win.DCBpDtrControl
	if on {
		dcb.Flags |= win.DTR_CONTROL_ENABLE << win.DCBpDtrControl
	} else {
		dcb.Flags |= win.DTR_CONTROL_DISABLE << win.DCBpDtrControl
	}
}

func (d *dev) updateCtl() (err error) {
	if d.inCtl {
		return
	}
	sav := &d.dcbsav
	dcb := &d.dcb
	if dcb.Flags == sav.Flags && dcb.BaudRate == sav.BaudRate &&
		dcb.ByteSize == sav.ByteSize &&
		dcb.Parity == sav.Parity &&
		dcb.StopBits == sav.StopBits {
		return
	}
	d.Drain()
	if e := win.SetCommState(d.fd, &d.dcb); e != nil {
		err = d.error("setdcb", e)
	} else {
		d.dcbsav = d.dcb
	}
	return
}

func (p *dev) SendBreak(ms int) error {
	return p.errorf("send-break", "not implemented")
}

func (d *dev) ModemLines() LineState {
	var ls LineState
	// TBD
	return ls
}
