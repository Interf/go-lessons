package structs

type Stack struct {
	doubleLink DoubleLink
}

func (s *Stack) Push(item int) {
	s.doubleLink.PushFront(item)
}

func (s *Stack) Pop() (int, bool) {

	node := s.doubleLink.PopFront()
	if node == nil {
		return 0, false
	}

	return node.Value, true
}
