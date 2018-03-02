package inmem

import "github.com/rafaeljesus/srv-consumer/pkg"

// Add a new user to the store.
func (s *Storage) Add(user *pkg.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, in := range s.users {
		if in.Username == user.Username {
			return pkg.ErrConflict
		}
	}

	user.ID = s.nextID(user)
	s.users[user.ID] = user

	return nil
}

// Save a user to the store.
func (s *Storage) Save(user *pkg.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[user.ID]; !ok {
		return pkg.ErrNotFound
	}

	s.users[user.ID] = user
	return nil
}
