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

func (s *Stack[T]) Peek() T {
	if s.length < 1 {
		var t T
		return t
	}

	return s.items[s.length-1]
}

func (s *Stack[T]) Peek2() T {
	if s.length < 2 {
		var t T
		return t
	}

	return s.items[s.length-2]
}

func (s *Stack[T]) Len() int {
	return s.length
}

func (s *Stack[T]) Clear() {
	s.items = nil
	s.length = 0
}
