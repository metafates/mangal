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
