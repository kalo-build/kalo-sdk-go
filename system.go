package kalo

import (
	"time"

	"github.com/kalo-build/kalo-sdk-go/hostabi"
)

// SystemCapabilities provides access to host system resources
// that are not available in standard WASI environments.
type SystemCapabilities struct{}

// System provides access to host system capabilities.
// Use this for resources like real-time clock that WASI doesn't expose.
var System = SystemCapabilities{}

// Now returns the current time from the host system.
// This bypasses WASI limitations where real-time clock access is restricted.
//
// Example usage:
//
//	timestamp := kalo.System.Now()
//	formatted := timestamp.Format("20060102150405")
func (s SystemCapabilities) Now() time.Time {
	nanos := hostabi.SystemNow()
	return time.Unix(0, nanos)
}

// NowUnix returns the current Unix timestamp in seconds from the host system.
func (s SystemCapabilities) NowUnix() int64 {
	return hostabi.SystemNow() / 1e9
}

// NowUnixNano returns the current Unix timestamp in nanoseconds from the host system.
func (s SystemCapabilities) NowUnixNano() int64 {
	return hostabi.SystemNow()
}
