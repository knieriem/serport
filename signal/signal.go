// +build darwin freebsd linux netbsd openbsd

// Package signal extends exp/signal by providing a dummy
// for the Incoming channel on systems that are not supported
// by exp/signal.
package signal

import (
	"os"
	"exp/signal"
)

var Incoming <-chan os.Signal

func init() {
	Incoming = signal.Incoming
}
