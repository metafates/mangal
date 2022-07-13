package util

import "sync"

// RwMap is a read-write map with support for concurrent reads and writes
// Generally faster than golang's sync.Map for reading / writing
type RwMap[K comparable, V any] struct {
	rw   sync.RWMutex
	data map[K]V
}

// Get returns the value for a key
func (m *RwMap[K, V]) Get(key K) (V, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	val, ok := m.data[key]
	return val, ok
}

// Set sets the value for a key
func (m *RwMap[K, V]) Set(key K, val V) {
	m.rw.Lock()
	defer m.rw.Unlock()

	m.data[key] = val
}

func NewRwMap[K comparable, V any]() *RwMap[K, V] {
	return &RwMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *RwMap[K, V]) Reset() {
	m.rw.Lock()
	defer m.rw.Unlock()

	m.data = make(map[K]V)
}

func (m *RwMap[K, V]) Len() int {
	m.rw.RLock()
	defer m.rw.RUnlock()

	return len(m.data)
}
