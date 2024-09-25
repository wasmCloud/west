// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package monotonicclock represents the imported interface "wasi:clocks/monotonic-clock@0.2.1".
//
// WASI Monotonic Clock is a clock API intended to let users measure elapsed
// time.
//
// It is intended to be portable at least between Unix-family platforms and
// Windows.
//
// A monotonic clock is a clock which has an unspecified initial value, and
// successive reads of the clock will produce non-decreasing values.
package monotonicclock

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	"go.wasmcloud.dev/wadge/bindings/wasi/io/poll"
)

// Instant represents the u64 "wasi:clocks/monotonic-clock@0.2.1#instant".
//
// An instant in time, in nanoseconds. An instant is relative to an
// unspecified initial value, and can only be compared to instances from
// the same monotonic-clock.
//
//	type instant = u64
type Instant uint64

// Duration represents the u64 "wasi:clocks/monotonic-clock@0.2.1#duration".
//
// A duration of time, in nanoseconds.
//
//	type duration = u64
type Duration uint64

// Now represents the imported function "now".
//
// Read the current value of the clock.
//
// The clock is monotonic, therefore calling this function repeatedly will
// produce a sequence of non-decreasing values.
//
//	now: func() -> instant
//
//go:nosplit
func Now() (result Instant) {
	result0 := wasmimport_Now()
	result = (Instant)((uint64)(result0))
	return
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.1 now
//go:noescape
func wasmimport_Now() (result0 uint64)

// Resolution represents the imported function "resolution".
//
// Query the resolution of the clock. Returns the duration of time
// corresponding to a clock tick.
//
//	resolution: func() -> duration
//
//go:nosplit
func Resolution() (result Duration) {
	result0 := wasmimport_Resolution()
	result = (Duration)((uint64)(result0))
	return
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.1 resolution
//go:noescape
func wasmimport_Resolution() (result0 uint64)

// SubscribeInstant represents the imported function "subscribe-instant".
//
// Create a `pollable` which will resolve once the specified instant
// has occurred.
//
//	subscribe-instant: func(when: instant) -> pollable
//
//go:nosplit
func SubscribeInstant(when Instant) (result poll.Pollable) {
	when0 := (uint64)(when)
	result0 := wasmimport_SubscribeInstant((uint64)(when0))
	result = cm.Reinterpret[poll.Pollable]((uint32)(result0))
	return
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.1 subscribe-instant
//go:noescape
func wasmimport_SubscribeInstant(when0 uint64) (result0 uint32)

// SubscribeDuration represents the imported function "subscribe-duration".
//
// Create a `pollable` that will resolve after the specified duration has
// elapsed from the time this function is invoked.
//
//	subscribe-duration: func(when: duration) -> pollable
//
//go:nosplit
func SubscribeDuration(when Duration) (result poll.Pollable) {
	when0 := (uint64)(when)
	result0 := wasmimport_SubscribeDuration((uint64)(when0))
	result = cm.Reinterpret[poll.Pollable]((uint32)(result0))
	return
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.1 subscribe-duration
//go:noescape
func wasmimport_SubscribeDuration(when0 uint64) (result0 uint32)
