// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// IsTerminal is based on a function in go.crypto/ssh/terminal/util.go

package ioutil

import (
	"syscall"
	"unsafe"
)

type FileDescriptor interface {
	Fd() uintptr
}

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(f FileDescriptor) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
