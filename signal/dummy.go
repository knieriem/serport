// +build !darwin,!freebsd,!linux,!netbsd,!openbsd

package signal

import (
	"os"
)

var Incoming = make(chan os.Signal, 0)
