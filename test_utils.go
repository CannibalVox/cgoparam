package cgoparam

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Test files can't call cgo so we have to put these as unexported non-test symbols
// compiler will optimize these out, so no worries

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

func callGoString(str unsafe.Pointer) string {
	return C.GoString((*C.char)(str))
}

func callGoBytes(b unsafe.Pointer, len int) []byte {
	return C.GoBytes(b, C.int(len))
}
