package structs

import "sync"

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

type DoubleLink struct {
	Head  *Node
	Tail  *Node
	Size  int
	mutex sync.Mutex
}

func (dl *DoubleLink) PushFront(value int) {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	node := &Node{
		Value: value,
	}

	if dl.Head == nil {
		dl.Head = node
		dl.Tail = node
		dl.Size++

		return
	}

	dl.Head.Prev = node
	node.Next = dl.Head
	dl.Head = node
	dl.Size++
}

func (dl *DoubleLink) PushBack(value int) {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	node := &Node{
		Value: value,
	}

	if dl.Head == nil {
		dl.Head = node
		dl.Tail = node
		dl.Size++

		return
	}

	dl.Tail.Next = node
	node.Prev = dl.Tail
	dl.Tail = node
	dl.Size++
}

func (dl *DoubleLink) PopFront() *Node {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	if dl.Size == 0 {
		return nil
	}

	node := dl.Head
	dl.Head = node.Next

	if node.Next != nil {
		node.Next.Prev = nil
	}

	if dl.Size == 1 {
		dl.Tail = node.Next
	}

	dl.Size--

	return node
}

func (dl *DoubleLink) PopBack() *Node {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	if dl.Size == 0 {
		return nil
	}

	node := dl.Tail
	dl.Tail = node.Prev

	if node.Prev != nil {
		node.Prev.Next = nil
	}

	if dl.Size == 1 {
		dl.Head = node.Prev
	}

	dl.Size--

	return node
}

func (dl *DoubleLink) FindByValue(value int) (*Node, bool) {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	node := dl.Head

	for i := 0; i < dl.Size; i++ {
		if value == node.Value {
			return node, true
		}

		node = node.Next
	}

	return nil, false
}

func (dl *DoubleLink) RemoveNodeByIndex(index int) *Node {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	headNode := dl.Head

	for i := 0; i < dl.Size; i++ {

		if index == i {

			if headNode.Prev == nil {
				dl.Head = headNode.Next

				if dl.Size == 1 {
					dl.Tail = headNode.Next
				}
				dl.Size--

				return headNode

			} else if headNode.Next == nil {
				dl.Tail = headNode.Prev

				if dl.Size == 1 {
					dl.Head = headNode.Prev
				}
				dl.Size--

				return headNode
			}

			headNode.Prev.Next = headNode.Next
			headNode.Next.Prev = headNode.Prev
			dl.Size--

			return headNode
		}

		headNode = headNode.Next
	}

	return nil
}

func (dl *DoubleLink) Len() int {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	return dl.Size
}

func (dl *DoubleLink) IsEmpty() bool {
	return dl.Len() == 0
}
