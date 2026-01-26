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

//go:wasmimport kalo db_get_migrations
func dbGetMigrations(storeID uint32) uint64

//go:wasmimport kalo db_apply_migration
func dbApplyMigration(storeID uint32, namePtr uint32, nameLen uint32, sqlPtr uint32, sqlLen uint32) uint32

//go:wasmimport kalo db_ensure_tracking_table
func dbEnsureTrackingTable(storeID uint32) uint32

// DBGetMigrations calls the host to get applied migrations for a store.
// Returns a pointer and length to a JSON-encoded array of AppliedMigration.
func DBGetMigrations(storeID uint32) (ptr uint32, length uint32) {
	packed := dbGetMigrations(storeID)
	return unpackPtrLen(packed)
}

// DBApplyMigration calls the host to apply a migration.
// Returns 0 on success, non-zero error code on failure.
func DBApplyMigration(storeID uint32, name string, sql []byte) uint32 {
	namePtr, nameLen := stringToPtr(name)
	sqlPtr, sqlLen := bytesToPtr(sql)
	return dbApplyMigration(storeID, namePtr, nameLen, sqlPtr, sqlLen)
}

// DBEnsureTrackingTable calls the host to ensure the migration tracking table exists.
// Returns 0 on success, non-zero error code on failure.
func DBEnsureTrackingTable(storeID uint32) uint32 {
	return dbEnsureTrackingTable(storeID)
}

// ReadMemory reads bytes from WASM linear memory at the given pointer.
func ReadMemory(ptr uint32, length uint32) []byte {
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

// unpackPtrLen unpacks a uint64 into (ptr, len).
// High 32 bits = pointer, low 32 bits = length.
func unpackPtrLen(packed uint64) (uint32, uint32) {
	ptr := uint32(packed >> 32)
	length := uint32(packed & 0xFFFFFFFF)
	return ptr, length
}

// stringToPtr converts a Go string to a pointer and length for passing to host functions.
func stringToPtr(s string) (uint32, uint32) {
	if len(s) == 0 {
		return 0, 0
	}
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

// bytesToPtr converts a byte slice to a pointer and length for passing to host functions.
func bytesToPtr(b []byte) (uint32, uint32) {
	if len(b) == 0 {
		return 0, 0
	}
	ptr := unsafe.Pointer(&b[0])
	return uint32(uintptr(ptr)), uint32(len(b))
}
