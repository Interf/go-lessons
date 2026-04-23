package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackPush(t *testing.T) {

	stack := Stack{}

	for i := 0; i < 5; i++ {
		stack.Push(i)
	}

	assert.Equal(t, 5, stack.doubleLink.Len())
}

func TestStackPop(t *testing.T) {
	stack := Stack{}

	for i := 1; i < 4; i++ {
		stack.Push(i)
	}

	tests := []struct {
		name string
		want int
	}{
		{
			name: "first item",
			want: 3,
		},
		{
			name: "second item",
			want: 2,
		},
		{
			name: "third item",
			want: 1,
		},
		{
			name: "empty item",
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ := stack.Pop()

			assert.Equal(t, test.want, got)
		})
	}

}
