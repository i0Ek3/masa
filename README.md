# masa

> masa: A memory allocation simulation algorithm.

masa is a simulation algorithm to allocate the memory. It supports level control allocation and status transfer. For now, masa rest on the experimental stage, we'll add new features soon.


## Level control allocation

- Tiny(`(0, 16)`)
- Little(`[16, 32]`)
- Enough(`(32, )`)


## Status transfer

- MemoryInitializing
- MemoryAllocated
- MemoryAllocating
- MemoryAllocateFailed
- MemoryWaitToAllocate


## License 

MIT.
