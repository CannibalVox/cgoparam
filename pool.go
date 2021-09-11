package cgoparam

import "sync"

const basePageSize = 4096
const considerStandaloneSize = 1024

var allocatorPool = sync.Pool{
	New: func() interface{} {
		return &Allocator{
			basePageSize: basePageSize,
			considerStandaloneSize: considerStandaloneSize,

			basePages: []*allocatorPage{createPage(basePageSize)},
		}
	},
}

func GetAlloc() *Allocator {
	return allocatorPool.Get().(*Allocator)
}

func ReturnAlloc(alloc *Allocator) {
	alloc.freeAll()
	allocatorPool.Put(alloc)
}
