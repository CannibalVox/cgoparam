package cgoparam

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

func allocAndDeallocSmallBatch() {
	xs := C.malloc(1)
	s := C.malloc(10)
	m := C.malloc(100)
	l := C.malloc(1000)
	C.free(l)
	C.free(m)
	C.free(s)
	C.free(xs)
}

func allocAndDeallocMultipage() {
	var allocs []unsafe.Pointer

	for i := 0; i < 10; i++ {
		allocs = append(allocs, C.malloc(1000))
	}

	for i := 0; i < 10; i++ {
		C.free(allocs[i])
	}
}
