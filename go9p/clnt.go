package go9p

import (
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
func Mount(net, addr, aname string, user p.User) (*Clnt, error) {
	c, err := clnt.Mount(net, addr, aname, user)
	return &Clnt{c}, err
}
