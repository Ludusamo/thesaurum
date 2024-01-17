package cache

import (
	"fmt"
)

type lruNode struct {
	next  *lruNode
	prev  *lruNode
	topic string
}

type lruTracker struct {
	head   *lruNode
	tail   *lruNode
	lookup map[string]*lruNode
}

func (n *lruNode) Print() {
	fmt.Printf("%s", n.topic)
	if n.next != nil {
		fmt.Print(" -> ")
		n.next.Print()
	}
}

func (n *lruNode) PrintReverse() {
	fmt.Printf("%s", n.topic)
	if n.prev != nil {
		fmt.Print(" -> ")
		n.prev.PrintReverse()
	}
}

func (tracker *lruTracker) Print() {
	if tracker.head != nil {
		fmt.Print("Forward: (")
		tracker.head.Print()
		fmt.Print(")\n")
	}
}

func (tracker *lruTracker) PrintReverse() {
	if tracker.tail != nil {
		fmt.Print("Reverse: (")
		tracker.tail.PrintReverse()
		fmt.Print(")\n")
	}
}

func newLruNode(topic string) *lruNode {
	var node lruNode
	node.next = nil
	node.prev = nil
	node.topic = topic
	return &node
}

func newLruTracker() *lruTracker {
	var tracker lruTracker
	tracker.head = nil
	tracker.tail = nil
	tracker.lookup = make(map[string]*lruNode)
	return &tracker
}

func (tracker *lruTracker) remove(topic string) *lruNode {
	node, found := tracker.lookup[topic]
	if found {
		prev := node.prev
		next := node.next
		if prev != nil {
			prev.next = next
		}
		if next != nil {
			next.prev = prev
		}
		if node == tracker.tail {
			tracker.tail = node.prev
			if node.prev != nil {
				node.prev.next = nil
			}
		}
		if node == tracker.head {
			tracker.head = node.next
			if node.next != nil {
				node.next.prev = nil
			}
		}
	}
	return node
}

func (tracker *lruTracker) Use(topic string) {
	fmt.Printf("using %s\n", topic)
	tracker.Print()
	tracker.PrintReverse()
	node, found := tracker.lookup[topic]
	if !found {
		fmt.Println("node not found, creating a new one")
		node = newLruNode(topic)
		tracker.lookup[topic] = node
	} else {
		tracker.remove(node.topic)
		node.prev = nil
		node.next = nil
	}

	head := tracker.head
	if head != nil {
		head.prev = node
		node.next = head
	}
	tracker.head = node

	if tracker.tail == nil {
		tracker.tail = node
	}
	fmt.Printf("after using %s\n", topic)
	tracker.Print()
	tracker.PrintReverse()
}

func (tracker *lruTracker) Pop() *lruNode {
	if tracker.tail != nil {
		lru := tracker.remove(tracker.tail.topic)
		if lru != nil {
			delete(tracker.lookup, lru.topic)
		}
		return lru
	}
	return nil
}

func (tracker *lruTracker) Delete(topic string) {
	n := tracker.remove(topic)
	if n != nil {
		delete(tracker.lookup, n.topic)
	}
}

type InMemoryStore struct {
	data       map[string]*Data
	nextStore  Store
	sizeCached int
	maxCached  int
	tracker    *lruTracker
}

func NewInMemoryStore(nextStore Store) *InMemoryStore {
	var s InMemoryStore
	s.data = make(map[string]*Data)
	s.nextStore = nextStore
	s.sizeCached = 0
	s.maxCached = 64
	s.tracker = newLruTracker()
	return &s
}

func (s *InMemoryStore) Store(topic string, data *Data) bool {
	fmt.Println("storing in memory")
	s.data[topic] = data
	s.sizeCached += data.Meta.Size
	for s.sizeCached > s.maxCached {
		fmt.Printf("%d is > %d maxCached\n", s.sizeCached, s.maxCached)
		node := s.tracker.Pop()
		lru, found := s.data[node.topic]
		if found {
			fmt.Printf("deleting %s from memory for %d\n", node.topic, lru.Meta.Size)
			s.Delete(node.topic)
		}
	}
	s.tracker.Use(topic)
	return true
}

func (s *InMemoryStore) Retrieve(topic string) (*Data, bool) {
	fmt.Println("retrieving in memory")
	data, found := s.data[topic]
	s.tracker.Use(topic)
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
