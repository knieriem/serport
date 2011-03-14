package main

import (
	"flag"
	"fmt"
	"log"
	"io"
	"os"
	"os/signal"
	"strings"
	"github.com/knieriem/g/sercom"
)

var (
	dev = flag.String("d", "", "COM device, e.g. COM1 or /dev/ttyUSB0")
	addr = flag.String("serve9P", "", "serve device via 9P at host:port")
	list = flag.Bool("list", false, "list serial devices")
	debug = flag.Bool("9d", false, "print 9P debug messages")
	debugall = flag.Bool("9D", false, "print 9P packets as well as debug messages")
)

func main() {
	var err os.Error
	var port sercom.Port

	flag.Parse()
	log.SetFlags(log.Lshortfile)

	sercom.Debug = *debug
	sercom.Debugall = *debugall

	if *list {
		for _, s := range sercom.DeviceList() {
			fmt.Println(s)
		}
		return
	}

	if strings.Index(*dev, ":") == -1 {
		port, err = sercom.Open(*dev, strings.Join(flag.Args(), " "))
	} else {
		port, err = sercom.Connect9P(*dev, "")
	}
	if err != nil {
		log.Fatalln(err)
	}
	if *addr != "" {
		go sercom.Serve9P(*addr, port)
	} else {
		go copyproc(port, os.Stdin)
		go copyproc(os.Stdout, port)
	}

	sig := <- signal.Incoming
	log.Println(sig)
	port.Close()
}

func copyproc(to io.Writer, from io.Reader) {
	var buf = make([]byte, 1024)

	for {
		n, err := from.Read(buf)
		if err != nil {
			if err != os.EOF {
				log.Fatalf("read: %s\n", err)
			}
			os.Exit(0)
		}
		if n > 0 {
			n, err = to.Write(buf[:n])
			if err != nil {
				log.Fatalf("write: %s\n", err)
			}
		}
	}
}
