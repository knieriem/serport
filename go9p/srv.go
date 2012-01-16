package go9p

import (
	"code.google.com/p/go9p/p"
	"syscall"
)

// Convert an error to a *p.Error
func ToError(err error) *p.Error {
	var ecode syscall.Errno

	if err == nil {
		return nil
	}

	ename := err.Error()
	if e, ok := err.(syscall.Errno); ok {
		ecode = e
	} else {
		ecode = syscall.EIO
	}

	return &p.Error{ename, ecode}
}
