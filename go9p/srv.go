package go9p

import (
	"os"
	"syscall"
	"go9p.googlecode.com/hg/p"
)

// Convert an os.Error to a *p.Error
func ToError(err os.Error) *p.Error {
	var ecode os.Errno

	if err == nil {
		return nil
	}

	ename := err.String()
	if e, ok := err.(os.Errno); ok {
		ecode = e
	} else {
		ecode = syscall.EIO
	}

	return &p.Error{ename, int(ecode)}
}
