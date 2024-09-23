// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package terminalstderr represents the imported interface "wasi:cli/terminal-stderr@0.2.1".
//
// An interface providing an optional `terminal-output` for stderr as a
// link-time authority.
package terminalstderr

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	terminaloutput "github.com/wasmCloud/west/examples/go/http/bindings/wasi/cli/terminal-output"
)

// GetTerminalStderr represents the imported function "get-terminal-stderr".
//
// If stderr is connected to a terminal, return a `terminal-output` handle
// allowing further interaction with it.
//
//	get-terminal-stderr: func() -> option<terminal-output>
//
//go:nosplit
func GetTerminalStderr() (result cm.Option[terminaloutput.TerminalOutput]) {
	wasmimport_GetTerminalStderr(&result)
	return
}

//go:wasmimport wasi:cli/terminal-stderr@0.2.1 get-terminal-stderr
//go:noescape
func wasmimport_GetTerminalStderr(result *cm.Option[terminaloutput.TerminalOutput])
