package mock

import "time"

type (
	Stats struct {
		StartInvoked bool
		TrackInvoked bool
		TrackFunc    func(t time.Time, err bool)
	}
)

func (s *Stats) Start() time.Time {
	s.StartInvoked = true
	return time.Now()
}

func (s *Stats) Track(t time.Time, err bool) {
	s.TrackInvoked = true
	s.TrackFunc(t, err)
}
