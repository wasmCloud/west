// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package exit represents the imported interface "wasi:cli/exit@0.2.0".
package exit

import (
	"github.com/ydnar/wasm-tools-go/cm"
)

// Exit represents the imported function "exit".
//
// Exit the current instance and any linked instances.
//
//	exit: func(status: result)
//
//go:nosplit
func Exit(status cm.BoolResult) {
	status0 := cm.BoolToU32(status)
	wasmimport_Exit((uint32)(status0))
	return
}

//go:wasmimport wasi:cli/exit@0.2.0 exit
//go:noescape
func wasmimport_Exit(status0 uint32)
