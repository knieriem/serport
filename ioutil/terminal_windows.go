// Copyright 2012 The g Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ioutil

import (
	"syscall"
)

type FileDescriptor interface {
	Fd() uintptr
}

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(f FileDescriptor) (is bool) {
	var mode uint32
	is = syscall.GetConsoleMode(syscall.Handle(f.Fd()), &mode) == nil
	return
}
