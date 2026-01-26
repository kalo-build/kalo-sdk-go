// Package kalo provides memory allocation exports for host function use.
package kalo

import (
	"sync"
	"unsafe"
)

// allocPool manages memory allocations for host function results.
// The host calls kalo_alloc to allocate memory, writes data, then
// the SDK reads and optionally frees the memory.
var allocPool = struct {
	sync.Mutex
	allocations map[uintptr][]byte
}{
	allocations: make(map[uintptr][]byte),
}

// kalo_alloc is exported to the host runtime.
// It allocates a byte slice of the given size and returns a pointer.
//
//go:export kalo_alloc
func kalo_alloc(size uint32) uint32 {
	if size == 0 {
		return 0
	}

	buf := make([]byte, size)
	ptr := uintptr(unsafe.Pointer(&buf[0]))

	allocPool.Lock()
	allocPool.allocations[ptr] = buf
	allocPool.Unlock()

	return uint32(ptr)
}

// kalo_free is exported to the host runtime.
// It frees a previously allocated pointer.
//
//go:export kalo_free
func kalo_free(ptr uint32) {
	if ptr == 0 {
		return
	}

	allocPool.Lock()
	delete(allocPool.allocations, uintptr(ptr))
	allocPool.Unlock()
}

// readAlloc reads the data at an allocated pointer.
// This is used internally by the SDK to read host function results.
func readAlloc(ptr, size uint32) []byte {
	if ptr == 0 || size == 0 {
		return nil
	}

	// Read directly from the pointer
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

