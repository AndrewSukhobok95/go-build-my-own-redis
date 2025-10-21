package storage

import (
	"log"
	"math"
	"path"
	"strconv"
	"sync"
	"time"
)

type KV struct {
	mu      sync.RWMutex
	data    map[string]*entry
	expires map[string]int64
}

func NewKV() *KV {
	return &KV{
		data:    make(map[string]*entry),
		expires: make(map[string]int64),
	}
}

func (s *KV) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = newStringEntry(value)
}

func (s *KV) Get(key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.data[key]
	if s.IsExpired(key) || !exists {
		return "", false, nil
	}

	strVal, err := e.String()
	return strVal, err == nil, err
}

func (s *KV) Delete(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int
	for _, k := range keys {
		if _, exists := s.data[k]; exists {
			delete(s.data, k)
			delete(s.expires, k)
			n++
		}
	}
	return n
}

func (s *KV) Incr(key string, delta int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var valInt int64 = 0
	e, exists := s.data[key]

	if exists {
		val, err := e.String()
		if err != nil {
			return 0, err
		}
		valInt, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, ErrNotInteger
		}
	}

	if (delta > 0 && valInt > math.MaxInt64-delta) || (delta < 0 && valInt < math.MinInt64-delta) {
		return 0, ErrOverflow
	}

	newValInt := valInt + delta
	s.data[key] = newStringEntry(strconv.FormatInt(newValInt, 10))
	return newValInt, nil
}

func (s *KV) Append(key, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var base string = ""
	if e, exists := s.data[key]; exists {
		valStr, err := e.String()
		if err != nil {
			return 0, err
		}
		base = valStr
	}

	newVal := base + value
	s.data[key] = newStringEntry(newVal)

	return len(newVal), nil
}

func (s *KV) Type(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, isp := s.data[key]
	if !isp || s.IsExpired(key) {
		return "none"
	}
	switch e.typ {
	case stringType:
		return "string"
	case listType:
		return "list"
	default:
		return "none"
	}
}

func (s *KV) Exists(keys ...string) int {
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

func (s *KV) Keys(pattern string) ([]string, error) {
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

func (s *KV) Flushdb() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.data)
	clear(s.expires)
}

func (s *KV) SetExpire(key string, duration time.Duration) (keyExists bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, keyExists = s.data[key]
	if keyExists {
		expireAt := time.Now().Add(duration)
		s.expires[key] = expireAt.UnixMilli()
	}
	return keyExists
}

func (s *KV) TTL(key string) int64 {
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

func (s *KV) IsExpired(key string) bool {
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

func (s *KV) ExpiredKeys() []string {
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

func (s *KV) Cleanup(interval time.Duration, stop <-chan struct{}) {
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

func (s *KV) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		s.data[key] = newListEntry(values)
		return len(values), nil
	}

	n, err := e.PushLeft(values...)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *KV) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		s.data[key] = newListEntry(values)
		return len(values), nil
	}

	n, err := e.PushRight(values...)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *KV) LPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		return "", nil
	}

	popped, err := e.PopLeft()
	if err != nil {
		return "", err
	}

	if len(e.data.([]string)) == 0 {
		delete(s.data, key)
	}
	return popped, nil
}

func (s *KV) RPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		return "", nil
	}

	popped, err := e.PopRight()
	if err != nil {
		return "", err
	}

	if len(e.data.([]string)) == 0 {
		delete(s.data, key)
	}
	return popped, nil
}

func (s *KV) LLen(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.data[key]
	if !exists {
		return 0, nil
	}

	n, err := e.LLen()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *KV) LRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.data[key]
	if !exists {
		return []string{}, nil
	}

	l, err := e.LRange(start, stop)
	if err != nil {
		return []string{}, err
	}
	return l, nil
}

func (s *KV) SAdd(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		e = newSetEntry()
		s.data[key] = e
	}
	cnt, err := e.SAdd(members...)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (s *KV) SMembers(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.data[key]
	if !exists {
		return []string{}, nil
	}

	members, err := e.SMembers()
	if err != nil {
		return []string{}, err
	}
	return members, nil
}

func (s *KV) SIsMember(key string, member string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.data[key]
	if !exists {
		return false, nil
	}

	isMember, err := e.SIsMember(member)
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func (s *KV) SRem(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		return 0, nil
	}

	cnt, err := e.SRem(members...)
	if err != nil {
		return 0, err
	}

	newSize, err := e.SLen()
	if err != nil {
		return 0, err
	}
	if newSize == 0 {
		delete(s.data, key)
	}
	return cnt, nil
}

func (s *KV) HSet(key, field, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		e = newHashEntry()
		s.data[key] = e
	}

	isNew, err := e.HSet(field, value)
	if err != nil {
		return 0, err
	}
	if isNew {
		return 1, nil
	}
	return 0, nil
}

func (s *KV) HGet(key, field string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.data[key]
	if !exists {
		return "", false, nil
	}

	value, fieldExists, err := e.HGet(field)
	if err != nil {
		return "", false, err
	}

	return value, fieldExists, nil
}

func (s *KV) HGetAll(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.data[key]
	if !exists {
		return []string{}, nil
	}

	m, err := e.HGetAll()
	if err != nil {
		return []string{}, err
	}

	flatHashSet := make([]string, 0, len(m)*2)
	for k, v := range m {
		flatHashSet = append(flatHashSet, k, v)
	}
	return flatHashSet, nil
}
