package kalo

import "errors"

var (
	// ErrStoreNotFound is returned when a requested store is not configured.
	ErrStoreNotFound = errors.New("store not found in configuration")

	// ErrInvalidStoreType is returned when a store is accessed with the wrong type.
	ErrInvalidStoreType = errors.New("invalid store type for requested operation")

	// ErrHostFunctionFailed is returned when a host function call fails.
	ErrHostFunctionFailed = errors.New("host function call failed")
)

