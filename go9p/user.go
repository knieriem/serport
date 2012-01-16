package go9p

import (
	"github.com/knieriem/g/syscall"
	"code.google.com/p/go9p/p"
	"os"
)

type user struct {
	name string
	id   int
}

func CurrentUser() p.User {
	return &user{syscall.GetUserName(), os.Getuid()}
}

func (u *user) Name() string {
	return u.name
}

func (u *user) Id() int {
	return u.id
}

func (u *user) Groups() []p.Group {
	return nil
}

func (u *user) IsMember(g p.Group) bool {
	return false
}
