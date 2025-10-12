package main

import (
	"log"
	"path"
	"sync"
	"time"
)

type Storage struct {
	mu      sync.RWMutex
	data    map[string]string
	expires map[string]int64
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, isp := s.data[key]
	if s.IsExpired(key) {
		return "", false
	}
	return val, isp
}

func (s *Storage) Delete(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int
	for _, k := range keys {
		if _, isp := s.data[k]; isp {
			delete(s.data, k)
			delete(s.expires, k)
			n++
		}
	}
	return n
}

func (s *Storage) Type(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, isp := s.data[key]
	if !isp || s.IsExpired(key) {
		return "none"
	}
	return "string"
}

func (s *Storage) Exists(keys ...string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var n int
	for _, k := range keys {
		if _, isp := s.data[k]; isp && !s.IsExpired(k) {
			n++
		}
	}
	return n
}

func (s *Storage) Keys(pattern string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var existing []string
	for k := range s.data {
		if s.IsExpired(k) {
			continue
		}
		matched, err := path.Match(pattern, k)
		if err != nil {
			log.Println("Error, while looking for pattern in keys:", err)
			return nil, err
		}
		if matched {
			existing = append(existing, k)
		}
	}
	return existing, nil
}

func (s *Storage) Flushdb() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.data)
	clear(s.expires)
}

func (s *Storage) SetExpire(key string, duration time.Duration) (keyExists bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, keyExists = s.data[key]
	if keyExists {
		expireAt := time.Now().Add(duration)
		s.expires[key] = expireAt.UnixMilli()
	}
	return keyExists
}

func (s *Storage) TTL(key string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, isp := s.data[key]
	if !isp {
		return -2
	}

	expiration, expired := s.expires[key]
	if !expired {
		return -1
	}

	ttl := time.Until(time.UnixMilli(expiration)).Milliseconds()
	if ttl < 0 {
		return -2
	}
	return ttl
}

func (s *Storage) IsExpired(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, isp := s.data[key]
	if !isp {
		return true
	}

	expiration, expired := s.expires[key]
	if !expired {
		return false
	}

	ttl := time.Until(time.UnixMilli(expiration))
	return ttl <= 0
}

func (s *Storage) ExpiredKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var expiredKeys []string
	for key := range s.expires {
		if s.IsExpired(key) {
			expiredKeys = append(expiredKeys, key)
		}
	}
	return expiredKeys
}

func (s *Storage) Cleanup(interval time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.Delete(s.ExpiredKeys()...)
		case <-stop:
			return
		}
	}
}
