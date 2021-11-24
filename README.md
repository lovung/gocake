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

