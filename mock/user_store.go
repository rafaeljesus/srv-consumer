package mock

import srv "github.com/rafaeljesus/srv-consumer"

type (
	UserStore struct {
		AddInvoked bool
		AddFunc    func(user *srv.User) error

		SaveInvoked bool
		SaveFunc    func(user *srv.User) error
	}
)

func (c *UserStore) Add(user *srv.User) error {
	c.AddInvoked = true
	return c.AddFunc(user)
}

func (c *UserStore) Save(user *srv.User) error {
	c.SaveInvoked = true
	return c.SaveFunc(user)
}
