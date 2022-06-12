package main

import "testing"

func TestContains(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	validItem := 7
	invalidItem := 42

	conditions := []bool{
		Contains[int](items, validItem),
		!Contains[int](items, invalidItem),
	}

	for _, condition := range conditions {
		if !condition {
			t.Fail()
		}
	}
}

func TestIfElse(t *testing.T) {
	IfElse(true, func() {}, func() { t.Fail() })()
	IfElse(false, func() { t.Fail() }, func() {})()
}

func TestPlural(t *testing.T) {
	conditions := []bool{
		Plural("word", 2) == "words",
		Plural("name", 29485) == "names",
		Plural("apple", 1) == "apple",
		Plural("dog", 0) == "dogs",
	}

	for _, condition := range conditions {
		if !condition {
			t.Fail()
		}
	}
}

func TestIsUnique(t *testing.T) {
	conditions := []bool{
		IsUnique([]int{1, 2, 3, 4}),
		!IsUnique([]int{1, 2, 3, 1}),
		IsUnique([]string{"Hello", "hello"}),
	}

	for _, condition := range conditions {
		if !condition {
			t.Fail()
		}
	}
}

func TestFind(t *testing.T) {
	integers := []int{1, 2, 3, 4, 10, 27, -258925}

	if found, ok := Find(integers, func(i int) bool {
		return i == 10
	}); ok {
		if found != 10 {
			t.Error("Wrong element was found")
		}
	} else {
		t.Error("Element was not found")
	}

	if _, ok := Find(integers, func(i int) bool {
		return i == 0
	}); ok {
		t.Error("ok is expected to be false for the non-existing element")
	}

	type person struct {
		name string
		age  int
	}

	structs := []person{
		{
			name: "name 1",
			age:  100,
		},
		{
			name: "name 2",
			age:  -1,
		},
	}

	if found, ok := Find(structs, func(p person) bool {
		return p.age < 0
	}); ok {
		if found.age != -1 {
			t.Error("Wrong element was found")
		}
	} else {
		t.Error("Element was not found")
	}
}
