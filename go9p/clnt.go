package go9p

import (
	"os"
	"go9p.googlecode.com/hg/p"
	"go9p.googlecode.com/hg/p/clnt"
)

const (
	OREAD   = p.OREAD
	OWRITE  = p.OWRITE
	ORDWR   = p.ORDWR
	OEXEC   = p.OEXEC
	OTRUNC  = p.OTRUNC
	OCEXEC  = p.OCEXEC
	ORCLOSE = p.ORCLOSE
)

type Clnt struct {
	*clnt.Clnt
}

// Same as go9p.googlecode.com/hg/p/clnt:Mount, but returning an os.Error
func Mount(net, addr, aname string, user p.User) (*Clnt, os.Error) {
	c, e9 := clnt.Mount(net, addr, aname, user)
	return &Clnt{c}, ToOsError(e9)
}


// convert a 9P error to an os.Error
func ToOsError(e9 *p.Error) os.Error {
	if e9 == nil {
		return nil
	}
	return os.NewError(e9.Error)
}
