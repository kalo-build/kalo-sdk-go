package kalo

import (
	"errors"
	"fmt"
)

var (
	// ErrStoreNotFound is returned when a requested store is not configured.
	ErrStoreNotFound = errors.New("store not found in configuration")

	// ErrInvalidStoreType is returned when a store is accessed with the wrong type.
	ErrInvalidStoreType = errors.New("invalid store type for requested operation")

	// ErrHostFunctionFailed is returned when a host function call fails.
	ErrHostFunctionFailed = errors.New("host function call failed")
)

// Host function error code meanings (must match kalo-cli/pkg/hostfuncs/db.go)
var hostErrorMessages = map[uint32]string{
	0: "success",
	1: "store not found",
	2: "failed to read name parameter",
	3: "failed to read SQL parameter",
	4: "failed to begin transaction",
	5: "migration SQL execution failed",
	6: "failed to record migration",
	7: "failed to commit transaction",
	8: "failed to create tracking table",
}

// HostFunctionError returns an error with the host function error code.
func HostFunctionError(operation string, errCode uint32) error {
	msg, ok := hostErrorMessages[errCode]
	if !ok {
		msg = "unknown error"
	}
	return fmt.Errorf("%s: %s (code %d)", operation, msg, errCode)
}

