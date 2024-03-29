package cache

import (
	"log"
	"strings"
	"sync"
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
	lock   *sync.Mutex
}

func (n *lruNode) constructString(sb *strings.Builder) {
	sb.WriteString(n.topic)
	if n.next != nil {
		sb.WriteString(" -> ")
		n.next.constructString(sb)
	}
}

func (n *lruNode) constructStringRev(sb *strings.Builder) {
	sb.WriteString(n.topic)
	if n.prev != nil {
		sb.WriteString(" -> ")
		n.prev.constructStringRev(sb)
	}
}

func (tracker *lruTracker) print() {
	var sb strings.Builder
	if tracker.head != nil {
		sb.WriteString("Forward: (")
		tracker.head.constructString(&sb)
		sb.WriteString(")\n")
		log.Print(sb.String())
	}
}

func (tracker *lruTracker) printReverse() {
	var sb strings.Builder
	if tracker.tail != nil {
		sb.WriteString("Reverse: (")
		tracker.tail.constructStringRev(&sb)
		sb.WriteString(")\n")
		log.Print(sb.String())
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
	tracker.lock = new(sync.Mutex)
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

func (tracker *lruTracker) use(topic string) {
	tracker.lock.Lock()
	defer tracker.lock.Unlock()
	node, found := tracker.lookup[topic]
	if !found {
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
}

func (tracker *lruTracker) pop() *lruNode {
	tracker.lock.Lock()
	defer tracker.lock.Unlock()
	log.Println("popping least recently used")
	tracker.print()
	if tracker.tail != nil {
		lru := tracker.remove(tracker.tail.topic)
		if lru != nil {
			delete(tracker.lookup, lru.topic)
		}
		return lru
	}
	return nil
}

func (tracker *lruTracker) delete(topic string) {
	tracker.lock.Lock()
	defer tracker.lock.Unlock()
	n := tracker.remove(topic)
	if n != nil {
		delete(tracker.lookup, n.topic)
	}
	tracker.print()
}
