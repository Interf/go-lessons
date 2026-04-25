package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	set := NewSet()

	set2 := &Set{
		data: make(map[int]struct{}),
	}

	assert.Equal(t, set, set2)
}

func TestAdd(t *testing.T) {
	set := NewSet()

	set.Add(1)

	_, ok := set.data[1]

	assert.True(t, ok)
}

func TestRemove(t *testing.T) {
	set := NewSet()
	set.Add(1)

	set.Remove(1)

	_, ok := set.data[1]

	assert.False(t, ok)
}

func TestContains(t *testing.T) {
	set := NewSet()
	set.Add(3)

	tests := map[string]struct {
		value int
		want  bool
	}{
		"exist element": {
			value: 3,
			want:  true,
		},
		"not exist element": {
			value: 1,
			want:  false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, set.Contains(tt.value))
		})
	}
}

func TestUnion(t *testing.T) {
	set1 := NewSet()

	for i := 0; i < 2; i++ {
		set1.Add(i)
	}

	set2 := NewSet()

	for i := 1; i < 5; i++ {
		set2.Add(i)
	}

	set3 := set1.Union(set2)

	for i := 0; i < 5; i++ {
		assert.True(t, set3.Contains(i))
	}
}

func TestIntersect(t *testing.T) {
	set1 := NewSet()

	for i := 0; i < 2; i++ {
		set1.Add(i)
	}

	set2 := NewSet()

	for i := 0; i < 5; i++ {
		set2.Add(i)
	}

	set3 := set1.Intersect(set2)

	for i := 0; i < 2; i++ {
		assert.True(t, set3.Contains(i))
	}
}

func TestDifference(t *testing.T) {
	set1 := NewSet()

	for i := 0; i < 2; i++ {
		set1.Add(i)
	}

	set2 := NewSet()

	for i := 0; i < 5; i++ {
		set2.Add(i)
	}

	set3 := set2.Difference(set1)

	for i := 2; i < 5; i++ {
		assert.True(t, set3.Contains(i))
	}
}
