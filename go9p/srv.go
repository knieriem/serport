package go9p

import (
	"code.google.com/p/go9p/p"
	"os"
	"syscall"
)

// Convert an os.Error to a *p.Error
func ToError(err error) *p.Error {
	var ecode os.Errno

	if err == nil {
		return nil
	}

	ename := err.Error()
	if e, ok := err.(os.Errno); ok {
		ecode = e
	} else {
		ecode = syscall.EIO
	}

	return &p.Error{ename, int(ecode)}
}
