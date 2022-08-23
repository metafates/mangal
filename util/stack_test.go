package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var stack = Stack[int]{}

func TestStack(t *testing.T) {
	Convey("Given an int stack", t, func() {
		Convey("When pushing a value", func() {
			stack.Push(1)
			Convey("Then the stack should have the value", func() {
				So(stack.Len(), ShouldEqual, 1)
				So(stack.Pop(), ShouldEqual, 1)
				So(stack.Len(), ShouldEqual, 0)
			})
		})

		Convey("When pushing multiple values", func() {
			stack.Push(1)
			stack.Push(2)
			stack.Push(3)
			Convey("Then the stack should have the values", func() {
				So(stack.Len(), ShouldEqual, 3)
				So(stack.Pop(), ShouldEqual, 3)
				So(stack.Pop(), ShouldEqual, 2)
				So(stack.Pop(), ShouldEqual, 1)
				So(stack.Len(), ShouldEqual, 0)
			})
		})
	})
}
