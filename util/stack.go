package util

type Stack[T any] struct {
	items  []*T
	length int
	last   *T
	first  *T
}

func (s *Stack[T]) Push(item *T) {
	s.items = append(s.items, item)
	s.length++
	s.last = s.items[s.length-1]
	s.first = s.items[0]
}

func (s *Stack[T]) Pop() *T {
	if s.length == 0 {
		var t *T
		return t
	}

	item := s.items[s.length-1]
	s.items = s.items[:s.length-1]
	s.length--
	s.last = s.items[s.length-1]
	s.first = s.items[0]
	return item
}

func (s *Stack[T]) Peek() *T {
	if s.length == 0 {
		var t *T
		return t
	}

	return s.items[s.length-1]
}

func (s *Stack[T]) Length() int {
	return s.length
}

func (s *Stack[T]) Clear() {
	s.items = nil
	s.length = 0
	s.last = nil
	s.first = nil
}

func (s *Stack[T]) Items() []*T {
	return s.items
}

func (s *Stack[T]) First() *T {
	return s.first
}
