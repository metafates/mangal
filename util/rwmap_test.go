package util

import "testing"

func TestNewRwMap(t *testing.T) {
	m := NewRwMap()

	if m == nil {
		t.Error("NewRwMap returned nil")
	}

	if len(m.data) != 0 {
		t.Error("NewRwMap returned non-empty map")
	}
}

func TestRwMap_Get(t *testing.T) {
	m := NewRwMap[string, string]()

	m.Set("test", "test")

	if v, ok := m.Get("test"); !ok {
		t.Error("Get returned false")
	} else if v != "test" {
		t.Error("Get returned invalid value")
	}

	if _, ok := m.Get("test2"); ok {
		t.Error("Get returned true")
	}
}

func TestRwMap_Set(t *testing.T) {
	m := NewRwMap[string, string]()

	m.Set("test", "test")

	if v, ok := m.Get("test"); !ok {
		t.Error("Get returned false")
	} else if v != "test" {
		t.Error("Get returned invalid value")
	}

	if _, ok := m.Get("test2"); ok {
		t.Error("Get returned true")
	}
}
