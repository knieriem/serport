package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"code.google.com/p/go9p/p"
	"code.google.com/p/go9p/p/srv"
	"github.com/knieriem/g/go9p"
	"github.com/knieriem/g/go9p/user"
	"github.com/knieriem/g/ioutil/terminal"
	"github.com/knieriem/g/netutil"
	"github.com/knieriem/serport"
	"github.com/knieriem/serport/encoding"
	"github.com/knieriem/serport/serenum"
	"github.com/knieriem/text/rc"
)

var (
	serveAddr = flag.String("serve", "", "serve device via 9P at a tcp `addr`, or (with `-') at stdin/out")

	list       = flag.Bool("list", false, "list serial devices")
	debug      = flag.Bool("9d", false, "print 9P debug messages")
	debugall   = flag.Bool("9D", false, "print 9P packets as well as debug messages")
	keepEcho   = flag.Bool("echo", false, "keep terminal's echo flag enabled")
	keepLine   = flag.Bool("line", false, "keep terminal's line flag enabled")
	crlfMode   = flag.Bool("crlf", false, "target needs CRLF line endings")
	binaryMode = flag.Bool("binary", false, "force binary mode (no modifications) even when using terminal")
	bridgePort = flag.String("bridge", "", "a serial `port` to copy data to and from")
)

type connOps struct {
	*srv.Fsrv
	*go9p.ConnOps
}

type traceLine struct {
	Δt     time.Duration
	prefix string
	buf    []byte
}

var traceC chan traceLine

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [DEVICE]\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	cherr = make(chan error)

	if *list {
		for _, info := range serenum.Ports() {
			fmt.Println(info.Device + "\n\t" + info.Format(nil))
		}
		return
	}

	devSpec := flag.Arg(0)
	if flag.NArg() > 0 {
		devSpec = strings.Join(flag.Args(), ",")
	}

	port, err := openport(devSpec)
	if err != nil {
		log.Fatalln(err)
	}

	restoreTerminal := func() {}

	if *serveAddr != "" {
		s, err := newServer(port)
		if err != nil {
			log.Fatal(err)
		}
		if *serveAddr != "-" {
			s.Start(s)
			go func() {
				cherr <- s.StartNetListener("tcp", *serveAddr)
			}()
		} else {
			o := go9p.AddConnOps(nil)
			closedC := o.ClosedC()
			s.Start(&connOps{s, o})
			s.NewConn(netutil.NewStreamConn(os.Stdin, os.Stdout))
			go func() {
				<-closedC
				cherr <- io.EOF
			}()
		}
	} else if *bridgePort != "" {
		port2, err := openport(*bridgePort)
		if err != nil {
			log.Fatal(err)
		}
		traceC = make(chan traceLine, 32)
		go func() {
			t0 := time.Now()
			for m := range traceC {
				ms := float64(m.Δt.Nanoseconds()) / 1000000
				t := time.Now()
				ms2 := float64(t.Sub(t0).Nanoseconds()) / 1000000
				if ms > 999 {
					fmt.Printf("-\t%06.2fms\t%s [%x]\n", ms2, m.prefix, m.buf)
				} else {
					fmt.Printf("%06.2fms\t%06.2fms\t%s [%x]\n", ms, ms2, m.prefix, m.buf)
				}
				t0 = t
			}
		}()
		go copyproc(port, port2, "<-")
		go copyproc(port2, port, "->")
	} else {
		if r := setupTerminal(os.Stdin); r != nil {
			restoreTerminal = r
		}

		// setup target encoding
		portRW := io.ReadWriter(port)
		switch {
		case *crlfMode:
			w := new(encoding.Wrapper)
			w.WrapReader(portRW, new(encoding.StripCR))
			w.WrapWriter(portRW, new(encoding.InsertCR))
			portRW = w
		default:
			w := new(encoding.Wrapper).WrapReader(portRW, new(encoding.StripCR))
			portRW = w.WrapWriter(portRW, nil)

		case *binaryMode:
		}

		// setup input encoding
		in := io.Reader(os.Stdin)
		if !*binaryMode {
			in = new(encoding.Wrapper).WrapReader(in, new(encoding.TermInput))
		}

		go copyproc(portRW, in, "")
		go copyproc(os.Stdout, portRW, "")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig)

	select {
	case err = <-cherr:
		if err != io.EOF {
			log.Println(err)
		}
	case s := <-sig:
		log.Println(s)
	}
	port.Close()
	restoreTerminal()
	os.Exit(0)
}

func setupTerminal(fd terminal.FileDescriptor) (restore func()) {
	if !terminal.IsTerminal(fd) {
		return
	}
	var clearFlags terminal.Flag
	if !*keepEcho {
		clearFlags |= terminal.EchoFlag
	}
	if !*keepLine {
		clearFlags |= terminal.LineFlag
	}
	oldState, err := terminal.DisableFlags(fd, clearFlags)
	if err != nil {
		log.Fatal(err)
	}
	restore = func() {
		terminal.Restore(fd, oldState)
	}
	return
}

const (
	defaultBaudrate = "b115200"
)

func openport(portSpec string) (port serport.Port, err error) {
	var c net.Conn

	f := strings.Split(portSpec, ",")
	dev := f[0]
	args := f[1:]

	for _, a := range args {
		if strings.HasPrefix(a, "b") {
			goto skipDefaultBaudrate
		}
	}

	args = append([]string{defaultBaudrate}, args...)
skipDefaultBaudrate:
	addr := dev
	details := ""

	prot := "local"
	dest := ""
	if i := strings.Index(dev, ":"); i != -1 {
		prot = dev[:i]
		dest = dev[i+1:]
	}
	if prot == "9P" {
		if strings.HasPrefix(dest, "!") {
			if c, err = connectToCommand(dest[1:]); err == nil {
				port, err = mountConn(c)
			}
		} else if c, err = net.Dial("tcp", dest); err == nil {
			port, err = mountConn(c)
		}
	} else if fi, e := os.Stat(dev); e == nil && fi.IsDir() {
		port, err = serport.OpenFsDev(dev)
	} else {
		var name string
		port, name, err = serport.Choose(dev, "")
		if err != nil {
			return
		}
		addr = name
		details = serenum.Lookup(name).Format(nil)
	}

	if len(args) != 0 {
		addr += "," + strings.Join(args, ",")
	}
	if err != nil {
		return
	}
	if len(args) != 0 {
		err = port.Ctl(strings.Join(args, " "))
		if err != nil {
			return
		}
	}
	fmt.Fprintln(os.Stderr, "# active device:", addr)
	if details != "" {
		fmt.Fprint(os.Stderr, "#\t(", details, ")\n")
	}
	return
}

func connectToCommand(command string) (c net.Conn, err error) {
	cmdLine := rc.Tokenize(command)

	cmd := exec.Command(cmdLine[0], cmdLine[1:]...)
	r, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	w, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	go func() {
		err := cmd.Wait()
		if err == nil {
			err = io.EOF
		}
		cherr <- err
	}()
	c = netutil.NewStreamConn(r, w)
	return
}

var cherr chan error

func copyproc(to io.Writer, from io.Reader, tracePrefix string) {
	var (
		buf = make([]byte, 1024)
		err error
		n   int
	)
	t0 := time.Now()

	for {
		if n, err = from.Read(buf); err != nil {
			break
		}
		if n > 0 {
			if tracePrefix != "" {
				t := time.Now()
				traceC <- traceLine{t.Sub(t0), tracePrefix, buf[:n]}
				t0 = t
			}
			if _, err = to.Write(buf[:n]); err != nil {
				break
			}
		}
	}
	cherr <- err
}

func mountConn(c net.Conn) (port serport.Port, err error) {
	port, clnt, err := serport.MountConn(c, "")
	if err == nil {
		switch {
		case *debugall:
			clnt.Debuglevel = 2
		case *debug:
			clnt.Debuglevel = 1

		}
	}
	return
}

func newServer(dev serport.Port) (s *srv.Fsrv, err error) {
	user := user.Current()
	root := new(srv.File)
	err = root.Add(nil, "/", user, nil, p.DMDIR|0555, nil)
	if err != nil {
		return
	}

	err = serport.RegisterFiles9P(root, dev, user)
	if err != nil {
		return
	}

	s = srv.NewFileSrv(root)

	switch {
	case *debugall:
		s.Debuglevel = 2
	case *debug:
		s.Debuglevel = 1
	}
	return
}
