package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueuePush(t *testing.T) {
	queue := Queue{}

	for i := 0; i < 5; i++ {
		queue.Push(i)
	}

	assert.Equal(t, 5, queue.doubleLink.Len())
}

func TestQueuePop(t *testing.T) {

	queue := Queue{}

	for i := 1; i < 4; i++ {
		queue.Push(i)
	}

	tests := []struct {
		name string
		want int
	}{
		{
			name: "first item",
			want: 1,
		},
		{
			name: "second item",
			want: 2,
		},
		{
			name: "third item",
			want: 3,
		},
		{
			name: "empty item",
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, _ := queue.Pop()

			assert.Equal(t, test.want, got)
		})
	}
}
