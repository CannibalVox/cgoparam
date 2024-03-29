package cgoparam

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

type allocatorPage struct {
	remainingSize int
	nextOffset    int
	size          int

	buffer unsafe.Pointer
}

func createPage(size int) *allocatorPage {
	ptr := C.malloc(C.size_t(size))
	return &allocatorPage{
		remainingSize: size,
		nextOffset:    0,
		size:          size,

		buffer: ptr,
	}
}

func (p *allocatorPage) Destroy() {
	C.free(p.buffer)
}

func (p *allocatorPage) FreeAll() {
	p.remainingSize = p.size
	p.nextOffset = 0
}

func (p *allocatorPage) NextPtr(size int) unsafe.Pointer {
	if p.remainingSize < size {
		panic("attempted to allocate more memory from page than it had. this indicates a disastrous bug in cgoparam")
	}

	ptr := unsafe.Add(p.buffer, p.nextOffset)
	oldOffset := p.nextOffset

	p.nextOffset += size
	alignment := p.nextOffset % 8
	if alignment != 0 {
		p.nextOffset = p.nextOffset - (p.nextOffset % 8) + 8
	}

	p.remainingSize -= p.nextOffset - oldOffset

	return ptr
}

type Allocator struct {
	basePageSize           int
	considerStandaloneSize int

	basePages        []*allocatorPage
	standaloneAllocs []unsafe.Pointer
}

func (a *Allocator) Malloc(size int) unsafe.Pointer {
	currentPage := a.basePages[len(a.basePages)-1]
	if size > currentPage.remainingSize {
		if size >= a.considerStandaloneSize {
			buffer := C.malloc(C.size_t(size))
			a.standaloneAllocs = append(a.standaloneAllocs, buffer)
			return buffer
		}

		newPage := createPage(a.basePageSize)
		a.basePages = append(a.basePages, newPage)
		return newPage.NextPtr(size)
	}

	return currentPage.NextPtr(size)
}

func (a *Allocator) CString(str string) unsafe.Pointer {
	strByteLen := len(str) + 1
	ptr := a.Malloc(strByteLen)
	ptrSlice := ([]byte)(unsafe.Slice((*byte)(ptr), strByteLen))
	copy(ptrSlice, str)
	ptrSlice[strByteLen-1] = 0
	return ptr
}

func (a *Allocator) CBytes(b []byte) unsafe.Pointer {
	byteLen := len(b)
	ptr := a.Malloc(byteLen)
	ptrSlice := ([]byte)(unsafe.Slice((*byte)(ptr), byteLen))
	copy(ptrSlice, b)
	return ptr
}

func (a *Allocator) freeAll() {
	basePageCount := len(a.basePages)
	standaloneCount := len(a.standaloneAllocs)

	a.basePages[0].FreeAll()
	for i := 1; i < basePageCount; i++ {
		a.basePages[i].Destroy()
	}
	a.basePages = a.basePages[:1]

	for i := 0; i < standaloneCount; i++ {
		C.free(a.standaloneAllocs[i])
	}
	a.standaloneAllocs = a.standaloneAllocs[:0]
}
