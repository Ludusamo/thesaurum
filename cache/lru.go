package cache

import "fmt"

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

func (n *lruNode) print() {
	fmt.Printf("%s", n.topic)
	if n.next != nil {
		fmt.Print(" -> ")
		n.next.print()
	}
}

func (n *lruNode) printReverse() {
	fmt.Printf("%s", n.topic)
	if n.prev != nil {
		fmt.Print(" -> ")
		n.prev.printReverse()
	}
}

func (tracker *lruTracker) print() {
	if tracker.head != nil {
		fmt.Print("Forward: (")
		tracker.head.print()
		fmt.Print(")\n")
	}
}

func (tracker *lruTracker) printReverse() {
	if tracker.tail != nil {
		fmt.Print("Reverse: (")
		tracker.tail.printReverse()
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

func (tracker *lruTracker) use(topic string) {
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
	fmt.Println("popping least recently used")
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
	n := tracker.remove(topic)
	if n != nil {
		delete(tracker.lookup, n.topic)
	}
}
