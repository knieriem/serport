// +build darwin freebsd linux netbsd openbsd

// Package signal extends exp/signal by providing a dummy
// for the Incoming channel on systems that are not supported
// by exp/signal.
package signal

import (
	"exp/signal"
	"os"
)

var Incoming <-chan os.Signal

func init() {
	Incoming = signal.Incoming
}
