package main

import (
	"sync"
	"time"
)

type idempotencyStore struct {
	mu      sync.Mutex
	expiry  map[string]time.Time
	maxKeys int
}

func newIdempotencyStore() *idempotencyStore {
	return &idempotencyStore{
		expiry:  make(map[string]time.Time),
		maxKeys: 10000,
	}
}

func (s *idempotencyStore) seenRecently(key string, ttl time.Duration) bool {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	if exp, ok := s.expiry[key]; ok {
		if now.Before(exp) {
			return true
		}
		delete(s.expiry, key)
	}

	if len(s.expiry) >= s.maxKeys {
		for k, exp := range s.expiry {
			if now.After(exp) {
				delete(s.expiry, k)
			}
		}
		if len(s.expiry) >= s.maxKeys {
			for k := range s.expiry {
				delete(s.expiry, k)
				break
			}
		}
	}

	s.expiry[key] = now.Add(ttl)
	return false
}
