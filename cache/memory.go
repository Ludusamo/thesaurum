package cache

import "fmt"

type InMemoryStore struct {
	data       map[string]*Data
	nextStore  Store
	sizeCached int
	maxCached  int
	tracker    *lruTracker
}

func NewInMemoryStore(nextStore Store, maxCached int) *InMemoryStore {
	var s InMemoryStore
	s.data = make(map[string]*Data)
	s.nextStore = nextStore
	s.sizeCached = 0
	s.maxCached = maxCached
	s.tracker = newLruTracker()
	return &s
}

func (s *InMemoryStore) Store(topic string, data *Data) bool {
	fmt.Println("storing in memory")
	s.data[topic] = data
	s.sizeCached += data.Meta.Size
	for s.sizeCached > s.maxCached {
		fmt.Printf("%d is > %d maxCached\n", s.sizeCached, s.maxCached)
		node := s.tracker.pop()
		lru, found := s.data[node.topic]
		if found {
			fmt.Printf("deleting %s from memory for %d\n", node.topic, lru.Meta.Size)
			s.Delete(node.topic)
		}
	}
	s.tracker.use(topic)
	return true
}

func (s *InMemoryStore) Retrieve(topic string) (*Data, bool) {
	fmt.Println("retrieving in memory")
	data, found := s.data[topic]
	s.tracker.use(topic)
	return data, found
}

func (s *InMemoryStore) Delete(topic string) bool {
	data, found := s.data[topic]
	if found {
		s.sizeCached -= data.Meta.Size
		delete(s.data, topic)
	}
	return true
}

func (s *InMemoryStore) List() []string {
	topics := make([]string, len(s.data))
	i := 0
	for k := range s.data {
		topics[i] = k
		i++
	}
	return topics
}

func (s *InMemoryStore) Next() Store {
	return s.nextStore
}

func (s *InMemoryStore) HasNext() bool {
	return s.nextStore != nil
}
