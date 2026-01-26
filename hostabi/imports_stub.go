//go:build !wasip1

// Package hostabi provides stub implementations for non-WASM builds.
// This allows the SDK to compile for native targets (testing, future transports).
package hostabi

// DBGetMigrations is a stub for non-WASM builds.
// In production, plugins are always compiled to WASM.
// This stub enables testing and future alternative transports (gRPC, etc).
func DBGetMigrations(storeID uint32) (ptr uint32, length uint32) {
	panic("hostabi: DBGetMigrations called in non-WASM build - use GOOS=wasip1 GOARCH=wasm")
}

// DBApplyMigration is a stub for non-WASM builds.
func DBApplyMigration(storeID uint32, name string, sql []byte) uint32 {
	panic("hostabi: DBApplyMigration called in non-WASM build - use GOOS=wasip1 GOARCH=wasm")
}

// DBEnsureTrackingTable is a stub for non-WASM builds.
func DBEnsureTrackingTable(storeID uint32) uint32 {
	panic("hostabi: DBEnsureTrackingTable called in non-WASM build - use GOOS=wasip1 GOARCH=wasm")
}

// ReadMemory is a stub for non-WASM builds.
func ReadMemory(ptr uint32, length uint32) []byte {
	panic("hostabi: ReadMemory called in non-WASM build - use GOOS=wasip1 GOARCH=wasm")
}
