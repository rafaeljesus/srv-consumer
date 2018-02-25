package inmem

import (
	"reflect"
	"sync"

	"github.com/rafaeljesus/srv-consumer/pkg"
)

type (
	// Storage manages in memory storage implementation.
	Storage struct {
		Driver  string
		mu      sync.RWMutex
		nextIDs map[interface{}]uint
		users   map[uint]*pkg.User
	}
)

// New creates a new in memory storage.
func New(dsn string) *Storage {
	return &Storage{
		Driver: "inmem",
	}
}

// nextID returns the next ID value that should be used for a struct of the given type.
func (s *Storage) nextID(val interface{}) uint {
	valType := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(val)).Interface())
	s.nextIDs[valType]++
	return s.nextIDs[valType]
}
