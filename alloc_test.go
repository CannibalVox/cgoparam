package cgoparam

import "testing"
import "github.com/stretchr/testify/require"

func TestBasicFunc(t *testing.T) {
	alloc := GetAlloc()
	require.NotNil(t, alloc)
	require.NotNil(t, alloc.Malloc(1))
	require.NotNil(t, alloc.Malloc(100))
	require.NotNil(t, alloc.Malloc(1000))

	require.Len(t, alloc.basePages, 1)
	require.Equal(t, 1101, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].size)
	require.Equal(t, 2995, alloc.basePages[0].remainingSize)

	ReturnAlloc(alloc)

	require.Len(t, alloc.basePages, 1)
	require.Equal(t, 0, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].size)
	require.Equal(t, 4096, alloc.basePages[0].remainingSize)
}

func TestMultiPage(t *testing.T) {
	alloc := GetAlloc()
	require.NotNil(t, alloc)
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))

	require.Len(t, alloc.basePages, 2)
	require.Len(t, alloc.standaloneAllocs, 0)
	require.Equal(t, 3600, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].size)
	require.Equal(t, 496, alloc.basePages[0].remainingSize)

	require.Equal(t, 900, alloc.basePages[1].nextOffset)
	require.Equal(t, 4096, alloc.basePages[1].size)
	require.Equal(t, 3196, alloc.basePages[1].remainingSize)

	ReturnAlloc(alloc)

	require.Len(t, alloc.basePages, 1)
	require.Equal(t, 0, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].remainingSize)
}

func TestSlightlyTooBigStandalone(t *testing.T) {
	alloc := GetAlloc()
	require.NotNil(t, alloc)
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(1200))

	require.Len(t, alloc.basePages, 1)
	require.Len(t, alloc.standaloneAllocs, 1)
	require.Equal(t, 3600, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].size)
	require.Equal(t, 496, alloc.basePages[0].remainingSize)

	ReturnAlloc(alloc)

	require.Len(t, alloc.basePages, 1)
	require.Len(t, alloc.standaloneAllocs, 0)
	require.Equal(t, 0, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].remainingSize)
}

func TestCouldBeStandaloneButFit(t *testing.T) {
	alloc := GetAlloc()
	require.NotNil(t, alloc)
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(900))
	require.NotNil(t, alloc.Malloc(1200))

	require.Len(t, alloc.basePages, 1)
	require.Len(t, alloc.standaloneAllocs, 0)
	require.Equal(t, 3900, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].size)
	require.Equal(t, 196, alloc.basePages[0].remainingSize)

	ReturnAlloc(alloc)

	require.Len(t, alloc.basePages, 1)
	require.Len(t, alloc.standaloneAllocs, 0)
	require.Equal(t, 0, alloc.basePages[0].nextOffset)
	require.Equal(t, 4096, alloc.basePages[0].remainingSize)
}

func TestCString(t *testing.T) {
	alloc := GetAlloc()
	defer ReturnAlloc(alloc)

	cStr := alloc.CString("WOW STRING")
	goStr := callGoString(cStr)
	require.Equal(t, "WOW STRING", goStr)
}

func TestCBytes(t *testing.T) {
	alloc := GetAlloc()
	defer ReturnAlloc(alloc)

	b := []byte("WOW STRING")

	cBytes := alloc.CBytes(b)
	goBytes := callGoBytes(cBytes, len(b))
	require.Equal(t, []byte("WOW STRING"), goBytes)
}
