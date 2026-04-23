package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleLinkPushFront(t *testing.T) {
	doubleLink := DoubleLink{}

	tests := map[string]struct {
		item int
		want int
	}{
		"first item": {
			item: 1,
			want: 1,
		},
		"second item": {
			item: 4,
			want: 4,
		},
		"third item": {
			item: 5,
			want: 5,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			doubleLink.PushFront(tt.item)

			assert.Equal(t, tt.want, doubleLink.Head.Value)
		})
	}

}

func TestDoubleLinkPushBack(t *testing.T) {
	doubleLink := DoubleLink{}

	tests := map[string]struct {
		item int
		want int
	}{
		"first item": {
			item: 5,
			want: 5,
		},
		"second item": {
			item: 1,
			want: 1,
		},
		"third item": {
			item: 2,
			want: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			doubleLink.PushBack(tt.item)

			assert.Equal(t, tt.want, doubleLink.Tail.Value)
		})
	}

}

func TestDoubleLinkPopFront(t *testing.T) {
	doubleLink := DoubleLink{}

	tests := map[string]struct {
		item int
		want *Node
	}{
		"first item": {
			item: 5,
			want: &Node{
				Value: 5,
			},
		},
		"second item": {
			item: 1,
			want: &Node{
				Value: 1,
			},
		},
		"third item": {
			item: 2,
			want: &Node{
				Value: 2,
			},
		},
		"empty item": {
			want: &Node{
				Value: 0,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.item > 0 {
				doubleLink.PushFront(tt.item)
			}

			if tt.item > 0 {
				assert.Equal(t, tt.want.Value, doubleLink.PopFront().Value)
			} else {
				assert.Nil(t, doubleLink.PopFront())
			}

		})
	}

}

func TestDoubleLinkPopBack(t *testing.T) {
	doubleLink := DoubleLink{}

	tests := map[string]struct {
		item int
		want *Node
	}{
		"first item": {
			item: 5,
			want: &Node{
				Value: 5,
			},
		},
		"second item": {
			item: 1,
			want: &Node{
				Value: 1,
			},
		},
		"third item": {
			item: 2,
			want: &Node{
				Value: 2,
			},
		},
		"empty item": {
			want: &Node{
				Value: 0,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.item > 0 {
				doubleLink.PushBack(tt.item)

				assert.Equal(t, tt.want.Value, doubleLink.PopBack().Value)
			} else {
				assert.Nil(t, doubleLink.PopBack())
			}

		})
	}
}

func TestDoubleLinkFindByValue(t *testing.T) {
	doubleLink := DoubleLink{}

	for i := 0; i < 5; i++ {
		doubleLink.PushBack(i)
	}

	tests := map[string]struct {
		item int
		want int
	}{
		"first item": {
			item: 4,
			want: 4,
		},
		"second item": {
			item: 1,
			want: 1,
		},
		"third item": {
			item: 2,
			want: 2,
		},
		"empty item": {
			item: 55,
			want: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			node, ok := doubleLink.FindByValue(tt.item)

			if ok {
				assert.Equal(t, tt.want, node.Value)
			} else {
				assert.Nil(t, node)
			}

		})
	}
}

func TestDoubleLinkRemoveNodeByIndex(t *testing.T) {

	doubleLink := DoubleLink{}

	for i := 1; i < 4; i++ {
		doubleLink.PushBack(i)
	}

	tests := []struct {
		name  string
		index int
		want  int
	}{
		{
			name:  "big item",
			index: 25,
			want:  0,
		},
		{
			name:  "first item",
			index: 0,
			want:  1,
		},
		{
			name:  "second item",
			index: 0,
			want:  2,
		},
		{
			name:  "third item",
			index: 3,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			node := doubleLink.RemoveNodeByIndex(tt.index)

			gotValue := 0

			if node != nil {
				gotValue = node.Value
			}

			assert.Equal(t, tt.want, gotValue)
		})
	}
}

func TestDoubleLinkLen(t *testing.T) {
	doubleLink := DoubleLink{}

	for i := 0; i < 5; i++ {
		doubleLink.PushBack(i)
	}

	assert.Equal(t, 5, doubleLink.Len())
}

func TestDoubleLinkIsEmpty(t *testing.T) {
	doubleLink := DoubleLink{}

	assert.True(t, doubleLink.IsEmpty())

	doubleLink.PushBack(1)
	assert.False(t, doubleLink.IsEmpty())
}
