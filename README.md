[![GoDoc](https://godoc.org/github.com/knieriem/serport?status.svg)](https://godoc.org/github.com/knieriem/serport)
[![GoReportCard](https://goreportcard.com/badge/github.com/knieriem/serport)](https://goreportcard.com/report/github.com/knieriem/serport)

Serport is a Go package providing access to serial ports on Linux
and Windows.  Its sub-package *serenum* helps finding serial ports on
a system.

In the following example a serial port is selected by `serport.Choose`;
if there is more than one serial port present on a system,
`Choose` will display a list of ports it has found and ask the user
to select one. If there is only one serial port present on a system,
`Choose` will try to use this port directly.
On success, the selected port is opened using 115200 baud, and obeying CTS.


	package main
	
	import (
		"fmt"
		"log"
	
		"github.com/knieriem/serport"
	)
	
	func main() {
		deviceName, err := serport.Choose("")
		if err != nil {
			log.Fatal(err)
		}
		conf := serport.MergeCtlCmds(serport.StdConf, "r1 m1")
		port, err := serport.Open(deviceName, conf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("# connected to", deviceName)

		// do something with `port'
		buf := make([]byte, 4096)
		n, err := port.Read(buf)
		fmt.Println(n, err)
	}

For a more complex example, see `cmd/sercom`.

For the syntax of control strings (for configuring baudrate, line control,
...) see [uart(3)](https://plan9.io/magic/man2html/3/uart).
