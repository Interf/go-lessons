package structs

type Item struct {
	Value int
}

type Stack struct {
	Items []Item
}

func (s *Stack) Push(i Item) {
	s.Items = append(s.Items, i)
}

func (s *Stack) Pop() Item {

	currentLen := len(s.Items)

	lastItem := s.Items[currentLen-1]

	newItems := make([]Item, currentLen-1)

	copy(newItems, s.Items[:currentLen-1])

	s.Items = newItems

	return lastItem
}
