package structs

type Queue struct {
	doubleLink DoubleLink
}

func (q *Queue) Push(item int) {
	q.doubleLink.PushBack(item)
}

func (q *Queue) Pop() (int, bool) {

	node := q.doubleLink.PopFront()
	if node == nil {
		return 0, false
	}

	return node.Value, true
}
