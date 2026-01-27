package kalo

import (
	"github.com/kalo-build/kalo-sdk-go/hostabi"
)

// DBStore provides generic database operations for a store.
// This is a thin transport layer - all business logic (migrations, etc.)
// should be implemented in the plugin that uses this interface.
type DBStore interface {
	// Exec executes SQL that doesn't return rows (INSERT, UPDATE, DELETE, DDL).
	Exec(sql []byte) error

	// Query executes SQL that returns rows and returns the result as JSON bytes.
	// The result format is an array of objects: [{"col1": val1, "col2": val2}, ...]
	Query(sql []byte) ([]byte, error)
}

// hostDBStore implements DBStore using host function callbacks.
type hostDBStore struct {
	storeName string
	storeID   uint32
	err       error
}

func (s *hostDBStore) Exec(sql []byte) error {
	if s.err != nil {
		return s.err
	}

	errCode := hostabi.DBExec(s.storeID, sql)
	if errCode != 0 {
		return HostFunctionError("exec sql", errCode)
	}

	return nil
}

func (s *hostDBStore) Query(sql []byte) ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}

	resultPtr, resultLen, errCode := hostabi.DBQuery(s.storeID, sql)
	if errCode != 0 {
		return nil, HostFunctionError("query sql", errCode)
	}

	if resultLen == 0 {
		return []byte("[]"), nil
	}

	return hostabi.ReadMemory(resultPtr, resultLen), nil
}
