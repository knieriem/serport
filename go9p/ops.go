package go9p

import (
	"code.google.com/p/go9p/p/srv"
)

func lazysetupChan(C *chan *srv.Conn) {
	if *C == nil {
		*C = make(chan *srv.Conn)
	}
}

func AddConnOps(s *srv.Srv) *ConnOps {
	co := new(ConnOps)
	co.Srv = s
	return co

}

type ConnOps struct {
	*srv.Srv
	closedC, openedC chan *srv.Conn
}

func (o *ConnOps) OpenedC() (c chan *srv.Conn) {
	lazysetupChan(&o.openedC)
	return o.openedC
}

func (o *ConnOps) ClosedC() (c chan *srv.Conn) {
	lazysetupChan(&o.closedC)
	return o.closedC
}

func (o *ConnOps) ConnOpened(c *srv.Conn) {
	if o.openedC != nil {
		o.openedC <- c
	}
}

func (o *ConnOps) ConnClosed(c *srv.Conn) {
	if o.closedC != nil {
		o.closedC <- c
	}
}
