package lru

import "sync"

type Node struct {
	Key   int
	Value int
	Prev  *Node
	Next  *Node
}

type LRUCache struct {
	data     map[int]*Node
	capacity int
	head     *Node
	tail     *Node
	mutex    sync.Mutex
}

const defaultCapacity = 8

func NewLRUCache(capacity int) *LRUCache {
	if capacity < 1 {
		capacity = defaultCapacity
	}

	head := &Node{}
	tail := &Node{}

	head.Next = tail
	tail.Prev = head

	return &LRUCache{
		data:     make(map[int]*Node, capacity),
		capacity: capacity,
		head:     head,
		tail:     tail,
	}

}

func (c *LRUCache) Get(key int) (int, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, ok := c.data[key]; ok {
		c.moveNode(node)

		return node.Value, true
	}

	return 0, false
}

func (c *LRUCache) Put(key int, value int) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, ok := c.data[key]; ok {
		node.Value = value
		c.moveNode(node)

		return
	}

	node := &Node{
		Key:   key,
		Value: value,
	}

	c.data[key] = node
	c.addToHead(node)

	if len(c.data) > c.capacity {
		c.removeTail()
	}
}

func (c *LRUCache) moveNode(node *Node) {
	c.removeNode(node)
	c.addToHead(node)
}

func (c *LRUCache) removeNode(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (c *LRUCache) addToHead(node *Node) {
	node.Prev = c.head
	node.Next = c.head.Next

	c.head.Next.Prev = node
	c.head.Next = node
}

func (c *LRUCache) removeTail() {
	tail := c.tail.Prev

	delete(c.data, tail.Key)
	c.removeNode(tail)
}
