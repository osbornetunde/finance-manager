package main

import (
	"container/list"
	"sync"
	"time"
)

type dedupItem struct {
	key string
	exp time.Time
}

type deduplicator struct {
	mu      sync.Mutex
	items   map[string]*list.Element
	evict   *list.List
	maxKeys int
}

func newDeduplicator() *deduplicator {
	return &deduplicator{
		items:   make(map[string]*list.Element),
		evict:   list.New(),
		maxKeys: 10000,
	}
}

func (s *deduplicator) seenRecently(key string, ttl time.Duration) bool {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	if ent, ok := s.items[key]; ok {
		if now.Before(ent.Value.(*dedupItem).exp) {
			s.evict.MoveToFront(ent)
			return true
		}
		// expired, remove it and let it fall through to be added fresh
		s.evict.Remove(ent)
		delete(s.items, key)
	}

	if s.evict.Len() >= s.maxKeys {
		ent := s.evict.Back()
		if ent != nil {
			s.evict.Remove(ent)
			delete(s.items, ent.Value.(*dedupItem).key)
		}
	}

	ent := s.evict.PushFront(&dedupItem{key: key, exp: now.Add(ttl)})
	s.items[key] = ent
	return false
}

func (s *deduplicator) forget(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ent, ok := s.items[key]
	if !ok {
		return
	}
	s.evict.Remove(ent)
	delete(s.items, key)
}
