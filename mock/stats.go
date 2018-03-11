package mock

import (
	"sync"
	"time"
)

type (
	Stats struct {
		sync.RWMutex
		StartInvoked bool
		TrackInvoked bool
		TrackFunc    func(t time.Time, err bool)
	}
)

func (s *Stats) Start() time.Time {
	s.Lock()
	defer s.Unlock()

	s.StartInvoked = true
	return time.Now()
}

func (s *Stats) Track(t time.Time, err bool) {
	s.Lock()
	defer s.Unlock()

	s.TrackInvoked = true
	s.TrackFunc(t, err)
}
