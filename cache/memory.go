package cache

import (
	"fmt"
	"sync"
)

type InMemoryStore struct {
	data       map[string]*Data
	sizeCached int
	maxCached  int
	tracker    *lruTracker
	lock       *sync.RWMutex
}

func NewInMemoryStore(maxCached int) *InMemoryStore {
	var s InMemoryStore
	s.data = make(map[string]*Data)
	s.sizeCached = 0
	s.maxCached = maxCached
	s.tracker = newLruTracker()
	s.lock = new(sync.RWMutex)
	return &s
}

func (s *InMemoryStore) put(topic string, data *Data) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[topic] = data
}

func (s *InMemoryStore) get(topic string) (*Data, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	d, found := s.data[topic]
	return d, found
}

func (s *InMemoryStore) Store(topic string, data *Data) error {
	fmt.Println("storing in memory")
	if data.Meta.Size > s.maxCached {
		// Skip caching since we would just wipe and still not be able to fit it
		return nil
	}
	s.put(topic, data)
	s.sizeCached += data.Meta.Size
	for s.sizeCached > s.maxCached {
		fmt.Printf("%d is > %d maxCached\n", s.sizeCached, s.maxCached)
		node := s.tracker.pop()
		lru, found := s.get(node.topic)
		if found {
			fmt.Printf("deleting %s from memory for %d\n", node.topic, lru.Meta.Size)
			s.Delete(node.topic)
		}
	}
	s.tracker.use(topic)
	return nil
}

func (s *InMemoryStore) Retrieve(topic string) (*Data, bool) {
	fmt.Println("retrieving in memory")
	data, found := s.get(topic)
	s.tracker.use(topic)
	return data, found
}

func (s *InMemoryStore) Delete(topic string) error {
	data, found := s.get(topic)
	if found {
		s.sizeCached -= data.Meta.Size
		s.lock.Lock()
		defer s.lock.Unlock()
		delete(s.data, topic)
	}
	return nil
}

func (s *InMemoryStore) List() []string {
	topics := make([]string, len(s.data))
	i := 0
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.data {
		topics[i] = k
		i++
	}
	return topics
}
