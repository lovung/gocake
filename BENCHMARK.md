# Current
```bash
❯ go test -bench=. -benchmem

goos: darwin
goarch: amd64
pkg: github.com/lovung/gocake
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkRistretto/Set-12         	1000000	     1018 ns/op	    518 B/op	      7 allocs/op
BenchmarkRistretto/Get-12         	2489055	      517.0 ns/op	     47 B/op	      3 allocs/op
BenchmarkGocake/Set-12            	1000000	     1096 ns/op	    221 B/op	      5 allocs/op
BenchmarkGocake/Get-12            	2796710	      385.7 ns/op	     33 B/op	      2 allocs/op
BenchmarkLFU/Touch-12             	3469186	      431.1 ns/op	    117 B/op	      0 allocs/op
BenchmarkLFU/Clean-12             	3404073	      662.0 ns/op	      8 B/op	      1 allocs/op
PASS
ok  	github.com/lovung/gocake	16.737s
```
# e6bbf0f
```bash
❯ go test -bench=. -benchmem

goos: darwin
goarch: amd64
pkg: github.com/lovung/gocake
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkRistretto/Set-12         	1255910	     1155 ns/op	    255 B/op	      7 allocs/op
BenchmarkRistretto/Get-12         	2706261	      413.5 ns/op	     47 B/op	      3 allocs/op
BenchmarkGocake/Set-12            	1000000	     1127 ns/op	    360 B/op	      5 allocs/op
BenchmarkGocake/Get-12            	2477674	      429.5 ns/op	     43 B/op	      2 allocs/op
BenchmarkLFU/Touch-12             	3530371	      417.3 ns/op	    115 B/op	      0 allocs/op
BenchmarkLFU/Clean-12             	3436413	      583.0 ns/op	      8 B/op	      1 allocs/op
PASS
ok  	github.com/lovung/gocake	15.536s
```
### 5d36142
```bash
❯ go test -bench=. -benchmem

goos: darwin
goarch: amd64
pkg: github.com/lovung/gocake
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkRistretto/Set-12         	1252419	     1169 ns/op	    244 B/op	      7 allocs/op
BenchmarkRistretto/Get-12         	2842155	      410.6 ns/op	     47 B/op	      3 allocs/op
BenchmarkGocake/Set-12            	1000000	     1132 ns/op	    360 B/op	      5 allocs/op
BenchmarkGocake/Get-12            	2555584	      421.4 ns/op	     43 B/op	      2 allocs/op
BenchmarkLFU/Touch-12             	3168571	      455.7 ns/op	    128 B/op	      0 allocs/op
BenchmarkLFU/Clean-12             	3438409	     5140 ns/op	      7 B/op	      0 allocs/op
PASS
ok  	github.com/lovung/gocake	31.322s
```

