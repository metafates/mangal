package stringutil

import (
	"golang.org/x/exp/constraints"
	"testing"
)

func TestFormatRanges(t *testing.T) {
	type args[T interface {
		constraints.Integer | constraints.Float
	}] struct {
		ranges []T
	}
	type testCase[T interface {
		constraints.Integer | constraints.Float
	}] struct {
		name string
		args args[T]
		want string
	}
	tests := []testCase[int]{
		{
			name: "empty",
			args: args[int]{},
			want: "",
		},
		{
			name: "single",
			args: args[int]{[]int{1}},
			want: "1",
		},
		{
			name: "double",
			args: args[int]{[]int{1, 2}},
			want: "1-2",
		},
		{
			name: "triple",
			args: args[int]{[]int{1, 2, 3}},
			want: "1-3",
		},
		{
			name: "semi-different",
			args: args[int]{[]int{1, 2, 4}},
			want: "1-2, 4",
		},
		{
			name: "different",
			args: args[int]{[]int{1, 3, 5}},
			want: "1, 3, 5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatRanges(tt.args.ranges); got != tt.want {
				t.Errorf("FormatRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuantify(t *testing.T) {
	type args struct {
		n        int
		singular string
		plural   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "zero",
			args: args{0, "manga", "mangas"},
			want: "0 mangas",
		},
		{
			name: "one",
			args: args{1, "manga", "mangas"},
			want: "1 manga",
		},
		{
			name: "two",
			args: args{2, "manga", "mangas"},
			want: "2 mangas",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quantify(tt.args.n, tt.args.singular, tt.args.plural); got != tt.want {
				t.Errorf("Quantify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		s   string
		max int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{"", 1},
			want: "",
		},
		{
			name: "short",
			args: args{"manga", 10},
			want: "manga",
		},
		{
			name: "long",
			args: args{"manga", 3},
			want: "ma…",
		},
		{
			name: "ansi",
			args: args{"\x1b[31mmanga\x1b[0m", 3},
			want: "\x1b[31mma…\x1b[0m",
		},
		{
			name: "unicode",
			args: args{"mangá", 3},
			want: "ma…",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.args.s, tt.args.max); got != tt.want {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}
