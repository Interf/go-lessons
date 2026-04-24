package lru

type Node struct {
	Key   int
	Value int
	Prev  *Node
	Next  *Node
}

type Cache struct {
	data     map[int]*Node
	head     *Node
	tail     *Node
	size     int
	capacity int
}

const DefaultCapacity = 8

func NewLRUCache(capacity int) *Cache {
	if capacity <= 0 {
		capacity = DefaultCapacity
	}

	return &Cache{
		data:     make(map[int]*Node, capacity),
		size:     0,
		capacity: capacity,
	}
}

func (c *Cache) Put(index int, value int) {
	if node, ok := c.data[index]; ok {
		node.Value = value
		c.removeNode(node)
		c.addToHead(node)

		return
	}

	node := &Node{
		Key:   index,
		Value: value,
	}

	c.data[index] = node
	c.addToHead(node)
	c.size++

	if c.size > c.capacity {
		c.removeTail()
	}

}

func (c *Cache) Get(index int) (int, bool) {
	node, ok := c.data[index]
	if !ok {
		return 0, false
	}

	if c.head != node {
		c.removeNode(node)
		c.addToHead(node)
	}

	return node.Value, true
}

func (c *Cache) removeNode(node *Node) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		c.head = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		c.tail = node.Prev
	}

	node.Next = nil
	node.Prev = nil
}

func (c *Cache) addToHead(node *Node) {
	node.Prev = nil
	node.Next = c.head

	if c.head != nil {
		c.head.Prev = node
	}

	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

func (c *Cache) removeTail() {
	if c.tail == nil {
		return
	}

	oldTail := c.tail

	c.removeNode(oldTail)
	delete(c.data, oldTail.Key)
	c.size--
}
