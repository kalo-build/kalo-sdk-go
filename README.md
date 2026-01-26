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
    GetAppliedMigrations() ([]AppliedMigration, error)
    ApplyMigration(name string, sql []byte) error
    EnsureTrackingTable() error
}
```

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
