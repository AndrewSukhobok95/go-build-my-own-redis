package main

import (
	"log"
	"path"
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]string
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
	return val, isp
}

func (s *Storage) Delete(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int
	for _, k := range keys {
		if _, ok := s.data[k]; ok {
			delete(s.data, k)
			n++
		}
	}
	return n
}

func (s *Storage) Type(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	if !ok {
		return "none"
	}
	return "string"
}

func (s *Storage) Exists(keys ...string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var n int
	for _, k := range keys {
		if _, ok := s.data[k]; ok {
			n++
		}
	}
	return n
}

func (s *Storage) Keys(pattern string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var existing []string
	for k, _ := range s.data {
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
}
