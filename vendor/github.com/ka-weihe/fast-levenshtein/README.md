# fast-levenshtein :rocket: 

> Fastest levenshtein implementation in Go.

> Measure the difference between two strings.

note: this implementation is currently not threadsafe and it assumes that the runes only go up to 65535. This will be fixed soon.

## Download
```bash
$ go get github.com/ka-weihe/fast-levenshtein
```

## Usage
```go
package main

import (
	"fmt"
	"github.com/ka-weihe/fast-levenshtein"
)

func main() {
	s1 := "fast"
	s2 := "fastest"
	distance := levenshtein.Distance(s1, s2)
	fmt.Printf("The distance between %s and %s is %d.\n", s1, s2, distance)
	// => The distance between fast and fastest is 3.
}
```

## Benchmarks
`kaweihe` is this package. It is 15 times faster for strings of length 64 compared to the second fastest package:

```bash
Benchmark/4/kaweihe-12         	   88651	     12858 ns/op	       0 B/op	       0 allocs/op
Benchmark/4/agniva-12          	   34550	     32099 ns/op	    7984 B/op	     499 allocs/op
Benchmark/4/arbovm-12          	   31778	     49311 ns/op	   23952 B/op	     499 allocs/op
Benchmark/4/dgryski-12         	   34489	     37951 ns/op	   23952 B/op	     499 allocs/op
Benchmark/8/kaweihe-12         	   51102	     22806 ns/op	       0 B/op	       0 allocs/op
Benchmark/8/agniva-12          	   17876	     66867 ns/op	   15968 B/op	     499 allocs/op
Benchmark/8/arbovm-12          	   12945	    104283 ns/op	   39920 B/op	     499 allocs/op
Benchmark/8/dgryski-12         	   11898	     92959 ns/op	   39920 B/op	     499 allocs/op
Benchmark/16/kaweihe-12        	   26749	     43723 ns/op	       0 B/op	       0 allocs/op
Benchmark/16/agniva-12         	    6129	    195921 ns/op	   23952 B/op	     499 allocs/op
Benchmark/16/arbovm-12         	    3370	    356006 ns/op	   71856 B/op	     499 allocs/op
Benchmark/16/dgryski-12        	    3242	    371193 ns/op	   71856 B/op	     499 allocs/op
Benchmark/32/kaweihe-12        	   12588	     93955 ns/op	       0 B/op	       0 allocs/op
Benchmark/32/agniva-12         	    1604	    767277 ns/op	   39920 B/op	     499 allocs/op
Benchmark/32/arbovm-12         	     838	   1434373 ns/op	  143712 B/op	     499 allocs/op
Benchmark/32/dgryski-12        	     837	   1427871 ns/op	  143712 B/op	     499 allocs/op
Benchmark/64/kaweihe-12        	    6538	    180181 ns/op	       0 B/op	       0 allocs/op
Benchmark/64/agniva-12         	     422	   2827182 ns/op	  327344 B/op	    1497 allocs/op
Benchmark/64/arbovm-12         	     219	   5449945 ns/op	  542912 B/op	    1497 allocs/op
Benchmark/64/dgryski-12        	     219	   5457779 ns/op	  542912 B/op	    1497 allocs/op
```

## License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
