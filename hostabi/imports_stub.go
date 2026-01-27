//go:build !wasip1

// Package hostabi provides stub implementations for non-WASM builds.
// This allows the SDK to compile for native targets (testing, future transports).
package hostabi

import "time"

// DBExec is a stub for non-WASM builds.
func DBExec(storeID uint32, sql []byte) uint32 {
	panic("hostabi: DBExec called in non-WASM build - use GOOS=wasip1 GOARCH=wasm")
}

// DBQuery is a stub for non-WASM builds.
func DBQuery(storeID uint32, sql []byte) (ptr uint32, length uint32, errCode uint32) {
	panic("hostabi: DBQuery called in non-WASM build - use GOOS=wasip1 GOARCH=wasm")
}

// ReadMemory is a stub for non-WASM builds.
func ReadMemory(ptr uint32, length uint32) []byte {
	panic("hostabi: ReadMemory called in non-WASM build - use GOOS=wasip1 GOARCH=wasm")
}

// SystemNow returns the current Unix timestamp in nanoseconds.
// For non-WASM builds, this uses the native time package.
func SystemNow() int64 {
	return time.Now().UnixNano()
}
