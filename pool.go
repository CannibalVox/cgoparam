package cgoparam

import (
	"runtime"
	"sync"
)

const basePageSize = 4096
const considerStandaloneSize = 1024

var allocatorPool = sync.Pool{
	New: func() interface{} {
		allocator := &Allocator{
			basePageSize:           basePageSize,
			considerStandaloneSize: considerStandaloneSize,

			basePages: []*allocatorPage{createPage(basePageSize)},
		}
		runtime.SetFinalizer(allocator, func(a *Allocator) {
			a.basePages[0].Destroy()
		})
		return allocator
	},
}

func GetAlloc() *Allocator {
	return allocatorPool.Get().(*Allocator)
}

func ReturnAlloc(alloc *Allocator) {
	alloc.freeAll()
	allocatorPool.Put(alloc)
}
