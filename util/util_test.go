package util

import (
	"golang.org/x/exp/slices"
	"regexp"
	"testing"
)

func TestContains(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	validItem := 7
	invalidItem := 42

	conditions := []bool{
		slices.Contains(items, validItem),
		!slices.Contains(items, invalidItem),
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

func TestMap(t *testing.T) {
	square := func(n int) int { return n * n }
	nums := []int{0, 1, 2, 3, 4}
	squared := Map(nums, square)

	if len(nums) != len(squared) {
		t.Error("Different lengths")
	}

	for i, s := range squared {
		n := nums[i]

		if square(n) != s {
			t.Error("Invalid value")
		}
	}
}

func TestToString(t *testing.T) {
	if ToString(1) != "1" {
		t.Error("Invalid value")
	}

	if ToString(1.0) != "1" {
		t.Error("Invalid value")
	}

	if ToString(true) != "true" {
		t.Error("Invalid value")
	}

	if ToString(false) != "false" {
		t.Error("Invalid value")
	}

	if ToString([]int{1, 2, 3}) != "[1 2 3]" {
		t.Error("Invalid value")
	}

	if ToString([]string{"a", "b", "c"}) != "[a b c]" {
		t.Error("Invalid value")
	}

	if ToString(map[string]int{"a": 1, "b": 2, "c": 3}) != "map[a:1 b:2 c:3]" {
		t.Error("Invalid value")
	}

	if ToString(struct{}{}) != "{}" {
		t.Error("Invalid value")
	}
}

func TestFetchLatestVersion(t *testing.T) {
	version, err := FetchLatestVersion()
	if err != nil {
		t.Error(err)
	}

	if version == "" {
		t.Error("Invalid version")
	}

	// make version regex
	versionRegex := regexp.MustCompile("^\\d+\\.\\d+\\.\\d+$")

	// check if version matches version regex
	if !versionRegex.MatchString(version) {
		t.Error("Invalid version")
	}

	// check if version is greater than 0.0.0
	if version < "0.0.0" {
		t.Error("Invalid version")
	}
}

func TestSanitizeFilename(t *testing.T) {

	// test with valid filename
	if SanitizeFilename("test.txt") != "test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test/test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test\\test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test:test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test*test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test?test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test|test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test<test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test>test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with invalid filename
	if SanitizeFilename("test?test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}

	// test with whitespace
	if SanitizeFilename("test test.txt") != "test_test.txt" {
		t.Error("Invalid filename")
	}
}

func TestPadZeros(t *testing.T) {
	if PadZeros(1, 2) != "01" {
		t.Error("Invalid value")
	}

	if PadZeros(10, 2) != "10" {
		t.Error("Invalid value")
	}

	if PadZeros(10, 10) != "0000000010" {
		t.Error("Invalid value")
	}

	if PadZeros(100, 2) != "100" {
		t.Error("Invalid value")
	}

}
