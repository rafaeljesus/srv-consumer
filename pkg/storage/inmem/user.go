package inmem

import (
	"errors"

	"github.com/rafaeljesus/srv-consumer/pkg"
)

func (s *Storage) Add(user *pkg.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, in := range s.users {
		if in.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	user.ID = s.nextID(user)
	s.users[user.ID] = user

	return nil
}

func (s *Storage) Save(user *pkg.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[user.ID]; !ok {
		return errors.New("user not found")
	}

	s.users[user.ID] = user
	return nil
}
