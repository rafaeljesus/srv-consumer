package mock

import (
	"github.com/rafaeljesus/srv-consumer/pkg"
)

type (
	UserStore struct {
		AddInvoked bool
		AddFunc    func(user *pkg.User) error

		SaveInvoked bool
		SaveFunc    func(user *pkg.User) error
	}
)

func (c *UserStore) Add(user *pkg.User) error {
	c.AddInvoked = true
	return c.AddFunc(user)
}

func (c *UserStore) Save(user *pkg.User) error {
	c.SaveInvoked = true
	return c.SaveFunc(user)
}
