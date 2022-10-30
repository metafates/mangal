package version

import (
	"fmt"
	"github.com/samber/lo"
	"strings"
)

func Compare(a, b string) (int, error) {
	type version struct {
		major, minor, patch int
	}

	parse := func(s string) (version, error) {
		var v version
		_, err := fmt.Sscanf(strings.TrimPrefix(s, "v"), "%d.%d.%d", &v.major, &v.minor, &v.patch)
		return v, err
	}

	av, err := parse(a)
	if err != nil {
		return 0, err
	}

	bv, err := parse(b)
	if err != nil {
		return 0, err
	}

	for _, pair := range []lo.Tuple2[int, int]{
		{av.major, bv.major},
		{av.minor, bv.minor},
		{av.patch, bv.patch},
	} {
		if pair.A > pair.B {
			return 1, nil
		}

		if pair.A < pair.B {
			return -1, nil
		}
	}

	return 0, nil
}
