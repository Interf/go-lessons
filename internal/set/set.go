package set

type Set struct {
	data map[int]struct{}
}

func NewSet() *Set {
	return &Set{
		data: make(map[int]struct{}),
	}
}

func (s *Set) Add(value int) {
	s.data[value] = struct{}{}
}

func (s *Set) Remove(value int) {
	delete(s.data, value)
}

func (s *Set) Contains(value int) bool {
	_, ok := s.data[value]
	return ok
}

func (s *Set) Union(other *Set) *Set {
	result := NewSet()

	for value := range s.data {
		result.Add(value)
	}

	for value, _ := range other.data {
		result.Add(value)
	}

	return result
}

func (s *Set) Intersect(other *Set) *Set {
	result := NewSet()

	for value := range s.data {
		if other.Contains(value) {
			result.Add(value)
		}
	}

	return result
}

func (s *Set) Difference(other *Set) *Set {
	result := NewSet()

	for value := range s.data {
		if !other.Contains(value) {
			result.Add(value)
		}
	}

	return result
}
