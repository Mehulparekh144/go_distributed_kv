package main

import (
	"fmt"
	"strconv"
	"time"
)

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]

	if !ok {
		return "", false
	}

	if value.TTL > 0 && time.Now().Unix() > value.TTL {
		return "", false
	}

	return value.value, true
}

func (s *Store) Set(key, value string, ttl int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logFile.WriteString(fmt.Sprintf("SET %s %s %s\n", key, value, strconv.FormatInt(ttl, 10)))

	expiration := int64(0)

	if ttl > 0 {
		expiration = time.Now().Unix() + ttl
	}

	s.data[key] = Item{
		value: value,
		TTL:   expiration,
	}
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logFile.WriteString(fmt.Sprintf("DELETE %s\n", key))

	delete(s.data, key)
}
