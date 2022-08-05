package util

type Stack[T any] struct {
	items  []T
	length int
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
	s.length++
}

func (s *Stack[T]) Pop() T {
	if s.length == 0 {
		var t T
		return t
	}

	item := s.items[s.length-1]
	s.items = s.items[:s.length-1]
	s.length--
	return item
}

func (s *Stack[T]) Length() int {
	return s.length
}

func (s *Stack[T]) Clear() {
	s.items = nil
	s.length = 0
}
