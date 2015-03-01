package serport

import (
	"os"

	"golang.org/x/sys/unix"

	sys "github.com/knieriem/g/syscall"
)

const (
	initDefault = "b115200 l8 pn s1"
)

type hw struct {
	*os.File
	inCtl          bool
	t, tsav, tOrig sys.Termios
}

// Open a local serial port.
// Inictl is an optional string containing one ore more commands in Ctl() style
func Open(filename string, inictl string) (port Port, err error) {
	var file *os.File

	// NONBLOCK prevents Open from blocking
	// until DCD is asserted from modem
	if file, err = os.OpenFile(filename, os.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0); err != nil {
		return
	}
	p := new(dev)
	p.File = file
	p.name = filename
	p.encaps = p
	t := &p.t

	fd := file.Fd()
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd, unix.TIOCEXCL, 0)
	if errno != 0 {
		err = errno
		goto fail
	}
	err = setBlocking(fd)
	if err != nil {
		err = p.error("set blocking", err)
		goto fail
	}

	if err = sys.IoctlTermios(fd, sys.TCGETS, t); err != nil {
		err = p.error("get term attr", err)
		goto fail
	}
	p.tOrig = p.t
	p.tsav = p.t

	t.Cflag |= sys.CLOCAL
	t.Lflag &^= sys.ICANON | sys.ECHO | sys.ISIG | sys.IEXTEN
	t.Iflag &^= sys.BRKINT | sys.ICRNL | sys.INPCK | sys.ISTRIP | sys.IXON
	t.Iflag |= sys.IGNPAR
	t.Oflag &^= sys.OPOST

	// block until at least one byte has been read
	t.Cc[sys.VMIN] = 1
	t.Cc[sys.VTIME] = 0

	if err = p.Ctl(initDefault, inictl); err != nil {
		goto fail
	}

	port = p
	return

fail:
	file.Close()
	return
}

func (p *dev) Close() error {
	sys.IoctlTermios(p.Fd(), sys.TCSETSW, &p.tOrig)
	return p.hw.Close()
}

func (d *dev) Drain() (err error) {
	err = sys.IoctlTermios(d.Fd(), sys.TCSETSW, &d.tsav) // drain and set parameters
	if err != nil {
		err = d.error("drain", err)
	}
	return
}

func (d *dev) Purge(in, out bool) {

}

func (d *dev) SetBaudrate(val int) (err error) {
	if speed, ok := speedMap[val]; !ok {
		err = d.errorf("open", "unsupported baud rate: %d", val)
		return
	} else {
		d.t.SetInSpeed(speed)
		d.t.SetOutSpeed(speed)
	}
	return d.updateCtl()
}

func (d *dev) SetWordlen(n int) error {
	t := &d.t

	t.Cflag &^= sys.CSIZE
	switch n {
	case 5:
		t.Cflag |= sys.CS5
	case 6:
		t.Cflag |= sys.CS6
	case 7:
		t.Cflag |= sys.CS7
	case 8:
		t.Cflag |= sys.CS8
	default:
		return d.errorf("open", "unsupported word len: %d", n)
	}
	return d.updateCtl()
}

func (d *dev) SetParity(parity byte) error {
	t := &d.t

	t.Cflag &^= sys.PARENB | sys.PARODD
	switch parity {
	case 'o':
		t.Cflag |= sys.PARODD
		fallthrough
	case 'e':
		t.Cflag |= sys.PARENB
	}
	return d.updateCtl()
}

func (d *dev) SetStopbits(n int) (err error) {
	switch n {
	case 1:
		d.t.Cflag &^= sys.CSTOPB
	case 2:
		d.t.Cflag |= sys.CSTOPB
	default:
		return d.errorf("open", "invalid number of stopbits: %d", n)
	}
	return d.updateCtl()
}

func (p *dev) SetRts(on bool) error {
	p.rts = on
	if on {
		return p.commfn("set rts", sys.TIOCMBIS, sys.TIOCM_RTS)
	}
	return p.commfn("clr rts", sys.TIOCMBIC, sys.TIOCM_RTS)
}
func (p *dev) SetDtr(on bool) error {
	p.dtr = on
	if on {
		return p.commfn("set dtr", sys.TIOCMBIS, sys.TIOCM_DTR)
	}
	return p.commfn("clr dtr", sys.TIOCMBIC, sys.TIOCM_DTR)
}

func (p *dev) commfn(name string, cmd int, f sys.Int) (err error) {
	if err = sys.IoctlModem(p.Fd(), cmd, &f); err != nil {
		return p.error(name, err)
	}
	return
}

func (d *dev) SetRtsCts(on bool) error {
	if on {
		d.t.Cflag |= sys.CRTSCTS
	} else {
		d.t.Cflag &^= sys.CRTSCTS
		d.SetRts(d.rts)
	}
	return d.updateCtl()
}

func (d *dev) updateCtl() (err error) {
	if d.inCtl {
		return
	}
	t := &d.t
	tsav := &d.tsav
	if t.Cflag == tsav.Cflag &&
		t.Lflag == tsav.Lflag &&
		t.Oflag == tsav.Oflag &&
		t.Cc[sys.VTIME] == tsav.Cc[sys.VTIME] &&
		t.Cc[sys.VMIN] == tsav.Cc[sys.VMIN] {
		return
	}
	if err = sys.IoctlTermios(d.Fd(), sys.TCSETSW, t); err == nil { // drain and set parameters
		d.tsav = d.t

		// It seems changing parameters also resets DTR/RTS lines;
		// put in the previously requested states again:
		d.SetRts(d.rts)
		d.SetDtr(d.dtr)
	}
	return
}

func (p *dev) ModemLines() LineState {
	var ls LineState
	return ls
}

var speedMap = map[int]int{
	50:      sys.B50,
	75:      sys.B75,
	110:     sys.B110,
	134:     sys.B134,
	150:     sys.B150,
	200:     sys.B200,
	300:     sys.B300,
	600:     sys.B600,
	1200:    sys.B1200,
	1800:    sys.B1800,
	2400:    sys.B2400,
	4800:    sys.B4800,
	9600:    sys.B9600,
	19200:   sys.B19200,
	38400:   sys.B38400,
	57600:   sys.B57600,
	115200:  sys.B115200,
	230400:  sys.B230400,
	460800:  sys.B460800,
	500000:  sys.B500000,
	576000:  sys.B576000,
	921600:  sys.B921600,
	1000000: sys.B1000000,
	1152000: sys.B1152000,
	1500000: sys.B1500000,
	2000000: sys.B2000000,
	2500000: sys.B2500000,
	3000000: sys.B3000000,
	3500000: sys.B3500000,
	4000000: sys.B4000000,
}

func setBlocking(fd uintptr) (err error) {
	var flags int

	flags, err = sys.Fcntl(fd, unix.F_GETFL, 0)
	if err == nil {
		_, err = sys.Fcntl(fd, unix.F_SETFL, flags&^unix.O_NONBLOCK)
	}
	return
}
