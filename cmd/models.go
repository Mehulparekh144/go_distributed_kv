package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Item struct {
	value string
	TTL   int64
}

type Store struct {
	data    map[string]Item
	mu      sync.RWMutex
	logFile *os.File
}

// NewStore creates a new Store with initialized data map
func NewStore(path string) *Store {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	store := &Store{
		data:    make(map[string]Item),
		logFile: f,
	}

	store.syncWAL()
	return store
}

func (s *Store) syncWAL() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logFile.Seek(0, 0)

	scanner := bufio.NewScanner(s.logFile)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "SET":
			if len(parts) == 4 {
				ttl, err := strconv.ParseInt(parts[3], 10, 64) // 10 is the base, 64 is the bit size
				if err != nil {
					panic(err)
				}
				s.data[parts[1]] = Item{
					value: parts[2],
					TTL:   ttl,
				}
			}
		case "DELETE":
			if len(parts) == 2 {
				delete(s.data, parts[1])
			}
		}
	}
}

func (s *Store) cleanupExpired() {
	go func() {
		for {
			fmt.Println("Cleaning up expired items")
			time.Sleep(1 * time.Second) // Sleep for 1 second
			now := time.Now().Unix()    // Get the current time

			s.mu.Lock()
			for key, item := range s.data {
				if item.TTL > 0 && now > item.TTL {
					fmt.Println("Deleting expired item:", key)
					delete(s.data, key)
					s.logFile.WriteString(fmt.Sprintf("DELETE %s\n", key))
				}
			}
			s.mu.Unlock()
		}
	}()
}

type SetRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   string `json:"ttl"` // Changed to string to handle JSON input
}

type HashRing struct {
	nodes        []string
	sortedHashes []uint32
	nodeMap      map[uint32]string
}

func NewHashRing(nodes []string) *HashRing {
	ring := &HashRing{
		nodes:   nodes,
		nodeMap: make(map[uint32]string),
	}

	for _, node := range nodes {
		h := crc32.ChecksumIEEE([]byte(node))
		ring.sortedHashes = append(ring.sortedHashes, h)
		ring.nodeMap[h] = node
	}

	sort.Slice(ring.sortedHashes, func(i, j int) bool { // Sort the hashes
		return ring.sortedHashes[i] < ring.sortedHashes[j]
	})

	return ring
}

func (r *HashRing) GetNode(key string) string {
	h := crc32.ChecksumIEEE([]byte(key))

	for _, hash := range r.sortedHashes {
		if hash >= h {
			return r.nodeMap[hash]
		}
	}

	return r.nodes[0]
}
