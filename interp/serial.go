package interp

import (
	"strings"

	"github.com/knieriem/serport"
	"github.com/knieriem/text/interp"
)

func NewCmdMap(port serport.Device, inictl []string) interp.CmdMap {
	return interp.CmdMap{
		"ctl": {
			Fn: func(_ interp.Context, arg []string) (err error) {
				return port.Ctl(arg[1:]...)
			},
			Arg:  []string{"CMD", "..."},
			Help: "Configure the serial port.",
		},
		"write": {
			Fn: func(_ interp.Context, arg []string) (err error) {
				b, err := interp.Argbytes(arg[1:])
				if err != nil {
					return
				}
				_, err = port.Write(b)
				return err
			},
			Arg:  []string{"BYTE", "..."},
			Help: "Write bytes to the serial port.",
		},
		"restore": {
			Fn: func(w interp.Context, arg []string) (err error) {
				if len(inictl) != 0 {
					w.Println("restoring", strings.Join(inictl, ","))
				}
				return port.Ctl(inictl...)
			},
			Help: "Restore serial port settings.",
		},
	}
}
