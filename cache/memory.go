package cache

import (
	"log"
	"sync"
)

type InMemoryCache struct {
	data       map[string]*Data
	sizeCached int
	maxCached  int
	tracker    *lruTracker
	lock       *sync.RWMutex
}

func NewInMemoryCache(maxCached int) *InMemoryCache {
	var s InMemoryCache
	s.data = make(map[string]*Data)
	s.sizeCached = 0
	s.maxCached = maxCached
	s.tracker = newLruTracker()
	s.lock = new(sync.RWMutex)
	return &s
}

func (s *InMemoryCache) put(topic string, data *Data) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[topic] = data
}

func (s *InMemoryCache) get(topic string) (*Data, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	d, found := s.data[topic]
	return d, found
}

func (s *InMemoryCache) store(topic string, data *Data) error {
	log.Println("storing in memory")
	if data.Meta.Size > s.maxCached {
		// Skip caching since we would just wipe and still not be able to fit it
		return nil
	}
	s.put(topic, data)
	s.sizeCached += data.Meta.Size
	for s.sizeCached > s.maxCached {
		log.Printf("%d is > %d maxCached\n", s.sizeCached, s.maxCached)
		node := s.tracker.pop()
		lru, found := s.get(node.topic)
		if found {
			log.Printf("deleting %s from memory for %d\n", node.topic, lru.Meta.Size)
			s.delete(node.topic)
		}
	}
	s.tracker.use(topic)
	return nil
}

func (s *InMemoryCache) retrieve(topic string) (*Data, bool) {
	log.Println("retrieving in memory")
	data, found := s.get(topic)
	if found {
		s.tracker.use(topic)
	}
	return data, found
}

func (s *InMemoryCache) delete(topic string) error {
	data, found := s.get(topic)
	if found {
		s.sizeCached -= data.Meta.Size
		s.lock.Lock()
		defer s.lock.Unlock()
		delete(s.data, topic)
		s.tracker.delete(topic)
	}
	return nil
}

func (s *InMemoryCache) list() []string {
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
