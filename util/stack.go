package util

// Stack is a LIFO data structure.
type Stack[T any] struct {
	items  []T
	length int
}

// Push adds an item to the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
	s.length++
}

// Pop removes an item from the stack.
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

// Peek returns the top item of the stack without removing it.
func (s *Stack[T]) Peek() T {
	if s.length < 1 {
		var t T
		return t
	}

	return s.items[s.length-1]
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	return s.length
}

// Clear removes all items from the stack.
func (s *Stack[T]) Clear() {
	s.items = nil
	s.length = 0
}
