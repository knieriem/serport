package serport

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	sys "github.com/knieriem/g/syscall"
)

type hw struct {
	*os.File
	t           *unix.Termios
	tsav, tOrig unix.Termios
	serOrig     *sys.Serial
	rs485       rs485State
	sc          syscall.RawConn

	rtsConfigured bool
	dtrConfigured bool
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
	d := new(dev)
	d.File = file
	d.name = filename
	d.encaps = d

	d.hw.sc, err = file.SyscallConn()
	if err != nil {
		return nil, err
	}

	var t *unix.Termios
	err = d.control(func(fd int) error {
		t, err = unix.IoctlGetTermios(fd, unix.TCGETS)
		if err != nil {
			err = d.error("get term attr", err)
			return err
		}
		return plainIoctl(fd, unix.TIOCEXCL)
	})
	if err != nil {
		goto fail
	}

	d.tOrig = *t
	d.tsav = *t
	d.t = t

	t.Cflag |= unix.CLOCAL
	t.Lflag &^= unix.ICANON | unix.ECHO | unix.ISIG | unix.IEXTEN
	t.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	t.Iflag |= unix.IGNPAR
	t.Oflag &^= unix.OPOST

	// block until at least one byte has been read
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0

	if err = d.Ctl(inictl); err != nil {
		goto fail
	}

	port = d
	return

fail:
	file.Close()
	return
}

func (d *dev) Read(buf []byte) (nread int, err error) {
	n, err := d.hw.Read(buf)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return 0, io.EOF
		}
		if errors.Is(err, os.ErrClosed) {
			return 0, io.EOF
		}
	}
	return n, err
}

func (d *dev) Close() error {
	d.hw.SetReadDeadline(time.Now().Add(-time.Second))
	d.control(func(fd int) error {
		unix.IoctlSetTermios(fd, unix.TCSETSW, &d.tOrig)
		d.rs485.restore(fd)
		if d.serOrig != nil {
			sys.IoctlSetSerial(fd, d.serOrig)
		}
		return nil
	})
	return d.hw.Close()
}

func (d *dev) Drain() error {
	return d.control(func(fd int) error {
		err := unix.IoctlSetTermios(fd, unix.TCSETSW, &d.tsav) // drain and set parameters
		if err != nil {
			return d.error("drain", err)
		}
		return nil
	})
}

func (d *dev) Purge(in, out bool) {

}

func (d *dev) SetBaudrate(val int) error {
	speed, ok := speedMap[val]
	if !ok {
		return d.errorf("open", "unsupported baud rate: %d", val)
	}
	c := d.t.Cflag &^ unix.CBAUD
	d.t.Cflag = c | (uint32(speed) & unix.CBAUD)
	return d.updateCtl()
}

func (d *dev) SetWordlen(n int) error {
	t := d.t

	t.Cflag &^= unix.CSIZE
	switch n {
	case 5:
		t.Cflag |= unix.CS5
	case 6:
		t.Cflag |= unix.CS6
	case 7:
		t.Cflag |= unix.CS7
	case 8:
		t.Cflag |= unix.CS8
	default:
		return d.errorf("open", "unsupported word len: %d", n)
	}
	return d.updateCtl()
}

func (d *dev) SetParity(parity byte) error {
	t := d.t

	t.Cflag &^= unix.PARENB | unix.PARODD
	switch parity {
	case 'o':
		t.Cflag |= unix.PARODD
		fallthrough
	case 'e':
		t.Cflag |= unix.PARENB
	}
	return d.updateCtl()
}

func (d *dev) SetStopbits(n int) (err error) {
	switch n {
	case 1:
		d.t.Cflag &^= unix.CSTOPB
	case 2:
		d.t.Cflag |= unix.CSTOPB
	default:
		return d.errorf("open", "invalid number of stopbits: %d", n)
	}
	return d.updateCtl()
}

func (d *dev) SetRts(on bool) error {
	d.rts = on
	d.rtsConfigured = true
	if on {
		return d.commfn("set rts", unix.TIOCMBIS, unix.TIOCM_RTS)
	}
	return d.commfn("clr rts", unix.TIOCMBIC, unix.TIOCM_RTS)
}
func (d *dev) SetDtr(on bool) error {
	d.dtr = on
	d.dtrConfigured = true
	if on {
		return d.commfn("set dtr", unix.TIOCMBIS, unix.TIOCM_DTR)
	}
	return d.commfn("clr dtr", unix.TIOCMBIC, unix.TIOCM_DTR)
}

func (d *dev) commfn(name string, cmd uint, f int) (err error) {
	return d.control(func(fd int) error {
		if err = unix.IoctlSetPointerInt(fd, cmd, f); err != nil {
			return d.error(name, err)
		}
		return nil
	})
}

func (d *dev) SetRtsCts(on bool) error {
	if on {
		d.t.Cflag |= unix.CRTSCTS
	} else {
		d.t.Cflag &^= unix.CRTSCTS
		d.SetRts(d.rts)
	}
	return d.updateCtl()
}

func (d *dev) updateCtl() (err error) {
	if d.inCtl {
		return
	}
	t := d.t
	tsav := &d.tsav
	if t.Cflag == tsav.Cflag &&
		t.Lflag == tsav.Lflag &&
		t.Oflag == tsav.Oflag &&
		t.Cc[unix.VTIME] == tsav.Cc[unix.VTIME] &&
		t.Cc[unix.VMIN] == tsav.Cc[unix.VMIN] {
		return
	}

	err = d.control(func(fd int) error {
		// drain and set parameters
		return unix.IoctlSetTermios(fd, unix.TCSETSW, t)
	})
	if err != nil {
		return err
	}

	d.tsav = *t

	// It seems changing parameters also resets DTR/RTS lines;
	// put in the previously requested states again:
	if d.rtsConfigured {
		d.SetRts(d.rts)
	}
	if d.dtrConfigured {
		d.SetDtr(d.dtr)
	}
	return
}

func (d *dev) SendBreak(duration time.Duration) error {
	if err := d.control(func(fd int) error {
		if err := plainIoctl(fd, unix.TIOCSBRK); err != nil {
			plainIoctl(fd, unix.TIOCCBRK)
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	time.Sleep(duration)
	return d.plainIoctl(unix.TIOCCBRK)
}

func (d *dev) ModemLines() LineState {
	var ls LineState
	return ls
}

var speedMap = map[int]int{
	50:      unix.B50,
	75:      unix.B75,
	110:     unix.B110,
	134:     unix.B134,
	150:     unix.B150,
	200:     unix.B200,
	300:     unix.B300,
	600:     unix.B600,
	1200:    unix.B1200,
	1800:    unix.B1800,
	2400:    unix.B2400,
	4800:    unix.B4800,
	9600:    unix.B9600,
	19200:   unix.B19200,
	38400:   unix.B38400,
	57600:   unix.B57600,
	115200:  unix.B115200,
	230400:  unix.B230400,
	460800:  unix.B460800,
	500000:  unix.B500000,
	576000:  unix.B576000,
	921600:  unix.B921600,
	1000000: unix.B1000000,
	1152000: unix.B1152000,
	1500000: unix.B1500000,
	2000000: unix.B2000000,
	2500000: unix.B2500000,
	3000000: unix.B3000000,
	3500000: unix.B3500000,
	4000000: unix.B4000000,
}

func (d *dev) plainIoctl(req uint) error {
	return d.control(func(fd int) error {
		return unix.IoctlSetInt(fd, req, 0)
	})
}

func plainIoctl(fd int, req uint) error {
	return unix.IoctlSetInt(fd, req, 0)
}

func (d *dev) SetLowLatency(low bool) error {
	return d.control(func(fd int) error {

		ser, err := sys.IoctlGetSerial(fd)
		if err != nil {
			return err
		}
		if d.serOrig == nil {
			orig := *ser
			d.serOrig = &orig
		}
		if low {
			ser.Flags |= sys.ASYNC_LOW_LATENCY
		} else {
			ser.Flags &^= sys.ASYNC_LOW_LATENCY
		}
		return sys.IoctlSetSerial(fd, ser)
	})
}

func init() {
	ctlNamespaceMap["rs485"] = rs485CtlNamespace
}

var rs485CtlNamespace = &ctlNamespace{
	runCmd: func(d *dev, cmd, c byte, n int) error {
		st := &d.rs485
		err := d.control(func(fd int) error {
			return st.initlazy(fd)
		})
		if err != nil {
			return fmt.Errorf("rs485: %w", err)
		}
		switch cmd {
		case 's':
			st.setFlag(sys.SER_RS485_RTS_ON_SEND, n)
		case 'a':
			st.setFlag(sys.SER_RS485_RTS_AFTER_SEND, n)
		case '[':
			st.cur.Rts_before_send = uint32(n)
		case ']':
			st.cur.Rts_after_send = uint32(n)
		case 't':
			st.setFlag(sys.SER_RS485_TERMINATE_BUS, n)
		case 'e':
			st.setFlag(sys.SER_RS485_RX_DURING_TX, n)
		default:
			return d.errorf("ctl", "rs485: unknown command: %q", string(cmd))
		}
		return nil
	},
	updateDrv: func(d *dev) error {
		err := d.control(func(fd int) error {
			return d.rs485.set(fd)
		})
		if err != nil {
			return fmt.Errorf("rs485: %w", err)
		}
		return nil
	},
}

func (d *dev) control(f func(fd int) error) error {
	var err error
	err1 := d.sc.Control(func(fd uintptr) {
		err = f(int(fd))
	})
	if err1 != nil {
		return err1
	}
	return err
}

type rs485State struct {
	cur  *sys.SerialRS485
	orig *sys.SerialRS485
}

func (st *rs485State) initlazy(fd int) error {
	if st.orig != nil {
		return nil
	}
	rs485, err := sys.IoctlGetSerialRS485(fd)
	if err != nil {
		return err
	}
	orig := *rs485
	st.orig = &orig
	st.cur = rs485
	rs485.Flags |= sys.SER_RS485_ENABLED
	return nil
}

func (st *rs485State) set(fd int) error {
	return sys.IoctlSetSerialRS485(fd, st.cur)
}

func (st *rs485State) setFlag(flag uint32, n int) {
	if n != 0 {
		st.cur.Flags |= flag
	} else {
		st.cur.Flags &^= flag
	}
}

func (st *rs485State) restore(fd int) error {
	if st.orig == nil {
		return nil
	}
	return sys.IoctlSetSerialRS485(fd, st.orig)
}
