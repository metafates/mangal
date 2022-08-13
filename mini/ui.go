package mini

import (
	"fmt"
	"github.com/metafates/mangal/style"
	"github.com/samber/lo"
	"os"
	"strconv"
	"strings"
)

func progress(msg string) (eraser func()) {
	msg = style.Combined(style.Blue, style.Truncate(truncateAt))(msg)
	fmt.Printf("\r%s", msg)

	return func() {
		_, _ = fmt.Fprintf(os.Stdout, "\r%s\r", strings.Repeat(" ", len(msg)))
	}
}

func title(t string) {
	fmt.Println(style.Combined(style.Magenta, style.Bold, style.Truncate(truncateAt))(t))
}

func fail(t string) {
	fmt.Println(style.Combined(style.Red, style.Bold, style.Truncate(truncateAt))(t))
}

func menu[T fmt.Stringer](items []T, options ...*bind) (*bind, T, error) {
	styles := map[int]func(string) string{
		0: style.Combined(style.Yellow, style.Truncate(truncateAt)),
		1: style.Combined(style.Cyan, style.Truncate(truncateAt)),
		2: style.Combined(style.Bold, style.Red, style.Truncate(truncateAt)),
	}

	for i, item := range items {
		s := fmt.Sprintf("(%d) %s", i+1, item.String())
		fmt.Println(styles[i%2](s))
	}

	options = append(options, quit)
	for i, option := range options {
		s := fmt.Sprintf("(%s) %s", option.A, option.B)
		s = style.Truncate(truncateAt)(s)

		if option == quit {
			fmt.Println(styles[2](s))
		} else {
			fmt.Println(styles[i%2](s))
		}
	}

	isValidOption := func(s string) bool {
		return lo.Contains(lo.Map(options, func(o *bind, _ int) string {
			return o.A
		}), s)
	}

	in, err := getInput(func(s string) bool {
		num, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return isValidOption(s)
		}
		return 0 < num && int(num-1) < len(items)+1
	})

	var t T

	if err != nil {
		return nil, t, err
	}

	if num, ok := in.asInt(); ok {
		return nil, items[num-1], nil
	}

	b, _ := parseBind(in.value)
	return b, t, nil
}
