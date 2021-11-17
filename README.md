# gocake
Golang in-memory cache library

## How it works
- [ ] Don't use the builtin `map` from Golang to avoid out of memory issue.
- [ ] Avoid to use pointer as much as possible to reduce GC overhead.
- [ ] Using CRFP (mix LRU/LFU) for eviction policy to optimize hit rate.
- [ ] Using WTinyLFU for admission policy to optimize hit rate.
- [ ] Implement own LRU and LFU based on our data structure.

## TODO
- [x] Cache interface definition.
- [ ] Basic cache feature.
- [ ] Sharding.
- [ ] Expired data clearing.
- [x] Thread safe.
- [ ] LRU.
- [x] LFU.
- [ ] CRFP.
