package structs

type Queue struct {
	Items []int
}

func (q *Queue) Push(item int) {
	q.Items = append(q.Items, item)
}

func (q *Queue) Pop() (int, bool) {

	currentLen := len(q.Items)

	if currentLen == 0 {
		return 0, false
	}

	firstItem := q.Items[0]

	newItems := make([]int, currentLen-1)

	copy(newItems, q.Items[1:currentLen])

	q.Items = newItems

	return firstItem, true
}
