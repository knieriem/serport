package serport

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/knieriem/g/ioutil/terminal"
	"github.com/knieriem/serport/serenum"
)

// Choose tries to get the name of a serial port by evaluating expr. On
// success it opens and returns the port, otherwise it returns an error.
// If expr does not equal one of the values "", "?" and "!", and does not
// start with "~", it is directly used as a port name. In the other cases, the
// system is queried for available ports. If necessary, and if os.Stdin is
// connected to a terminal, the user is prompted to select a port from
// the list, or to enter a port name if the list is empty.
func Choose(expr string, inictl string) (port Port, name string, err error) {
	switch expr {
	case "", "?", "!":
		name, err = choosePort(expr)
		if err != nil {
			return
		}
	default:
		if strings.HasPrefix(expr, "~") {
			name, err = matchDevice(expr[1:])
			if err != nil {
				return
			}
		} else {
			name = expr
		}
	}
	port, err = Open(name, inictl)
	return
}

func matchDevice(expr string) (name string, err error) {
	rx, err := regexp.Compile(expr)
	if err != nil {
		return
	}
	list := serenum.Ports()
	if len(list) == 0 {
		err = errors.New("no devices found")
		return
	}
	for _, p := range list {
		s := fmt.Sprintf("%v (%v)", p.Device, p.Format(nil))
		if rx.MatchString(s) {
			name = p.Device
			return
		}
	}
	err = errors.New("no matching device")
	return
}

func choosePort(mode string) (name string, err error) {
	sep := " .. "
	list := serenum.Ports()
	switch len(list) {
	case 0:
	case 1:
		sep = ""
		if mode != "?" {
			name = list[0].Device
			return
		}
	case 2:
		sep = ", "
		fallthrough
	default:
		if mode == "!" {
			name = list[0].Device
			return
		}
	}

	if !terminal.IsTerminal(os.Stdin) {
		return name, errors.New("os.Stdin is not connected to a terminal")
	}
	t, err := terminal.OpenOutput()
	if err != nil {
		return
	}
	defer t.Close()

	if len(list) == 0 {
		fmt.Fprint(t, "Enter serial port: ")
		_, err = fmt.Scan(&name)
		return
	}
	fmt.Fprintln(t, "\nChoose a serial port: ")
	for i, p := range list {
		fmt.Fprintf(t, "  %d\t%v (%v)\n", i, p.Device, p.Format(nil))
	}
	if sep == "" {
		fmt.Fprint(t, "\nPress return to select the device.")
	} else {
		fmt.Fprint(t, "\nEnter a number {(0)", sep, len(list)-1, "}: ")
	}

	s := bufio.NewScanner(os.Stdin)
	if !s.Scan() {
		err = s.Err()
		return
	}
	input := strings.TrimSpace(s.Text())
	var i int
	if input != "" {
		i64, err1 := strconv.ParseInt(input, 10, 32)
		if err1 != nil {
			err = err1
			return
		}
		i = int(i64)
	}
	if i >= len(list) {
		err = errors.New("value exceeds maximum index")
		return
	}
	name = list[i].Device
	return
}
