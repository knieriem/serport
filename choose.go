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

// Choose looks up serial ports available on a system,
// and selects one port depending on the value of expr.
// On success it returns the name of the port device found,
// which can be used in a subsequent call to [Open];
// otherwise it returns an error.
//
// When evaluating expr, the following rules apply:
//
//   - expr equals "":
//     The system is queried for available ports.
//     If only one port is found, its name is returned.
//     Otherwise the user is prompted to select from a list,
//     or, if no port could be found, to enter a port name.
//
//   - expr equals "?":
//     Like "", except that also in case only one port could be found,
//     the user is prompted to select from a list.
//
//   - expr equals "!":
//     Like "", except that if more than one port have been found,
//     the name of the first one is returned; no prompt is displayed in that case.
//
//   - expr starts with '~':
//     The system is queried for available ports.
//     If one or more ports have been found, the
//     part of expr following the ~ is used as a regular expression
//     on port names, device description, and serial number.
//     The name of the first matching port is returned.
//     In case there is no match, an error is returned.
//
//   - default: expr is returned unchanged
//
// Note that, in the above cases, if the user needs to be prompted,
// this is done only if [os.Stdin] is connected to a terminal;
// otherwise an error is returned.
func Choose(expr string) (name string, err error) {
	switch expr {
	case "", "?", "!":
		return choosePort(expr)
	default:
		if strings.HasPrefix(expr, "~") {
			return matchDevice(expr[1:])
		}
	}
	return expr, nil
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
