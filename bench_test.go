package cgoparam

import "testing"

func BenchmarkRawCgoSmallBatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allocAndDeallocSmallBatch()
	}
}

func cgoparamSmallBatch() {
	alloc := GetAlloc()
	defer ReturnAlloc(alloc)

	_ = alloc.Malloc(1)
	_ = alloc.Malloc(10)
	_ = alloc.Malloc(100)
	_ = alloc.Malloc(1000)
}

func BenchmarkCgoparamSmallBatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cgoparamSmallBatch()
	}
}

func BenchmarkRawCgoMultipage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allocAndDeallocMultipage()
	}
}

func cgoparamMultipage() {
	alloc := GetAlloc()
	defer ReturnAlloc(alloc)

	for i := 0; i < 10; i++ {
		_ = alloc.Malloc(1000)
	}
}

func BenchmarkCgoparamMultipage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cgoparamMultipage()
	}
}
