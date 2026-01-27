# Kalo SDK for Go

The Go-language SDK for Kalo CLI plugin development. This SDK provides a clean abstraction layer for plugins to access stores configured in `kalo.yaml`, regardless of whether they're backed by local filesystems, databases, or other network resources.

## Installation

```bash
go get github.com/kalo-build/kalo-sdk-go
```

## Usage

Plugin authors use the SDK to access stores by name:

```go
package main

import (
    "fmt"
    "github.com/kalo-build/kalo-sdk-go"
)

func main() {
    // Access a filesystem store
    fs := kalo.FS("KA_MIGRATIONS")
    
    // Read files
    data, err := fs.ReadFile("001_initial.sql")
    if err != nil {
        panic(err)
    }
    
    // List directory contents
    entries, err := fs.ListDir(".")
    if err != nil {
        panic(err)
    }
    
    // Access a database store
    db := kalo.DB("DB_MAIN")
    
    // Get applied migrations
    applied, err := db.GetAppliedMigrations()
    if err != nil {
        panic(err)
    }
    
    // Apply a migration
    err = db.ApplyMigration("002_add_users.sql", sqlContent)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Migration complete!")
}
```

## System Capabilities

The SDK provides access to host system resources not available in standard WASI:

```go
import kalo "github.com/kalo-build/kalo-sdk-go"

// Get current time (bypasses WASI clock restrictions)
now := kalo.System.Now()           // time.Time
unix := kalo.System.NowUnix()      // int64 (seconds)
nanos := kalo.System.NowUnixNano() // int64 (nanoseconds)

// Format for migration timestamps
timestamp := kalo.System.Now().Format("20060102150405")
```

This is essential for plugins that need real-time clock access (e.g., generating timestamped migration files).

## Store Types

The SDK automatically handles different store types:

### `localFileSystem`
- Accessed via `kalo.FS(storeName)`
- Uses WASI filesystem operations
- Mounted by CLI at runtime

### `cloudSqlDatabase`
- Accessed via `kalo.DB(storeName)`
- Uses host function callbacks
- CLI handles actual database connections

## Interfaces

### FileStore

```go
type FileStore interface {
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, data []byte) error
    ListDir(path string) ([]DirEntry, error)
    Stat(path string) (FileInfo, error)
    Exists(path string) bool
}
```

### DBStore

```go
type DBStore interface {
    // Exec executes SQL that doesn't return rows (INSERT, UPDATE, DELETE, DDL).
    Exec(sql []byte) error

    // Query executes SQL that returns rows and returns the result as JSON bytes.
    // The result format is an array of objects: [{"col1": val1, "col2": val2}, ...]
    Query(sql []byte) ([]byte, error)
}
```

The SDK provides only generic database operations. All business logic (migrations,
tracking, etc.) should be implemented in the plugin that uses this interface.

## Building WASM Plugins

Plugins using this SDK must be compiled for WASI:

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm ./cmd/plugin
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Plugin Code                          │
│  fs := kalo.FS("STORE")    db := kalo.DB("STORE")       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Kalo SDK                              │
│  FileStore (WASI)         DBStore (Host Functions)      │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Kalo CLI                              │
│  WASI FS Mounts          Host Function Registry         │
│  (/input, /output)       (db_get_migrations, etc.)      │
└─────────────────────────────────────────────────────────┘
```

## License

MIT License
