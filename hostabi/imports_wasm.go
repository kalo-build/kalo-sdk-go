//go:build wasip1

// Package hostabi provides the low-level host function imports for WASM plugins.
// Plugin authors should NOT use this package directly - use the kalo package instead.
package hostabi

import "unsafe"

// Host function imports for database operations.
// These are provided by the Kalo CLI runtime.
//
// ABI Note: go:wasmimport only allows 1 return value.
// For functions returning (ptr, len), we pack them into a uint64:
//   - High 32 bits: pointer
//   - Low 32 bits: length
//
// For functions returning (ptr, len, errCode), we use two uint64 values
// but go:wasmimport only allows one return, so we use a packed format:
//   - Return uint64 for error code (0 = success)
//   - Result data is written to a callback buffer

//go:wasmimport kalo db_exec
func dbExec(storeID uint32, sqlPtr uint32, sqlLen uint32) uint32

//go:wasmimport kalo db_query
func dbQuery(storeID uint32, sqlPtr uint32, sqlLen uint32) uint64

// Host function imports for system capabilities.
// These provide access to host system resources not available in WASI.

//go:wasmimport kalo system_now
func systemNow() int64

// DBExec calls the host to execute SQL that doesn't return rows.
// Returns 0 on success, non-zero error code on failure.
func DBExec(storeID uint32, sql []byte) uint32 {
	sqlPtr, sqlLen := bytesToPtr(sql)
	return dbExec(storeID, sqlPtr, sqlLen)
}

// DBQuery calls the host to execute SQL that returns rows.
// Returns (ptr, len, errCode) where ptr/len point to JSON result data.
// Error code is in low 8 bits, ptr in high 32 bits of first return,
// length in second return.
func DBQuery(storeID uint32, sql []byte) (ptr uint32, length uint32, errCode uint32) {
	sqlPtr, sqlLen := bytesToPtr(sql)
	packed := dbQuery(storeID, sqlPtr, sqlLen)
	// Packed format: high 32 bits = ptr, low 32 bits = (len << 8 | errCode)
	ptr = uint32(packed >> 32)
	lenAndErr := uint32(packed & 0xFFFFFFFF)
	length = lenAndErr >> 8
	errCode = lenAndErr & 0xFF
	return
}

// ReadMemory reads bytes from WASM linear memory at the given pointer.
func ReadMemory(ptr uint32, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

// bytesToPtr converts a byte slice to a pointer and length for passing to host functions.
func bytesToPtr(b []byte) (uint32, uint32) {
	if len(b) == 0 {
		return 0, 0
	}
	ptr := unsafe.Pointer(&b[0])
	return uint32(uintptr(ptr)), uint32(len(b))
}

// SystemNow calls the host to get the current Unix timestamp in nanoseconds.
// This provides real-time clock access since WASI environments typically don't have it.
func SystemNow() int64 {
	return systemNow()
}
