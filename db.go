package kalo

import (
	"encoding/json"

	"github.com/kalo-build/kalo-sdk-go/hostabi"
)

// DBStore provides database operations for a store.
type DBStore interface {
	// GetAppliedMigrations returns the list of migrations that have been applied.
	GetAppliedMigrations() ([]AppliedMigration, error)

	// ApplyMigration applies a migration with the given name and SQL content.
	ApplyMigration(name string, sql []byte) error

	// EnsureTrackingTable ensures the migration tracking table exists.
	EnsureTrackingTable() error
}

// AppliedMigration represents a migration that has been applied to the database.
type AppliedMigration struct {
	Name      string `json:"name"`
	Checksum  string `json:"checksum"`
	AppliedAt int64  `json:"appliedAt"` // Unix timestamp
}

// hostDBStore implements DBStore using host function callbacks.
type hostDBStore struct {
	storeName string
	storeID   uint32
	err       error
}

func (s *hostDBStore) GetAppliedMigrations() ([]AppliedMigration, error) {
	if s.err != nil {
		return nil, s.err
	}

	// Call host function to get migrations
	resultPtr, resultLen := hostabi.DBGetMigrations(s.storeID)
	if resultLen == 0 {
		return []AppliedMigration{}, nil
	}

	// Read result from WASM memory
	resultBytes := hostabi.ReadMemory(resultPtr, resultLen)

	var migrations []AppliedMigration
	if err := json.Unmarshal(resultBytes, &migrations); err != nil {
		return nil, err
	}

	return migrations, nil
}

func (s *hostDBStore) ApplyMigration(name string, sql []byte) error {
	if s.err != nil {
		return s.err
	}

	// Call host function to apply migration
	errCode := hostabi.DBApplyMigration(s.storeID, name, sql)
	if errCode != 0 {
		return ErrHostFunctionFailed
	}

	return nil
}

func (s *hostDBStore) EnsureTrackingTable() error {
	if s.err != nil {
		return s.err
	}

	errCode := hostabi.DBEnsureTrackingTable(s.storeID)
	if errCode != 0 {
		return ErrHostFunctionFailed
	}

	return nil
}

