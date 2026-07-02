package memory

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found in cache")

type entry struct {
	value     []byte
	expiresAt time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]entry),
	}
}

func (s *Store) set(key string, value []byte, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = entry{value: value, expiresAt: time.Now().Add(ttl)}
}

func (s *Store) get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, ErrNotFound
	}
	return e.value, nil
}

func (s *Store) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (s *Store) deleteByPrefix(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range s.data {
		if strings.HasPrefix(k, prefix) {
			delete(s.data, k)
		}
	}
}
