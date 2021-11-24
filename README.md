# gocake
Golang in-memory cache library

## How it work
- [x] Sharding the hashed key for each stores.
- [ ] Using LFU for eviction policy.
- [x] Thread safe.
- [x] Handle expired async.

## Future
- [ ] Don't use the builtin `map` from Golang to avoid out of memory issue.
- [ ] Avoid to use pointer as much as possible to reduce GC overhead.
- [ ] Using CRFP (mix LRU/LFU) for eviction policy to optimize hit rate.
- [ ] Using WTinyLFU for admission policy to optimize hit rate.

## Benchmark
### Current
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
