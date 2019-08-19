package etherr

import "fmt"

type RpcError struct {
	err     string
	Code    int
	Details string
}

func (e *RpcError) Error() string {
	return fmt.Sprintf("%s %s", e.err, e.Details)
}

// Nil is returned when a rpc call returned a null result
var Nil = New("Result is null", 0, "")

// ConnectionClosed is returned when the websocket connection is closed
var ConnectionClosed = New("Websocket connection closed", 0, "")

// Empty is returned when a rpc call returned an empty result
var Empty = New("Result is empty", 0, "")

// InvalidUInt8 is returned when a rpc call returned an empty result
var InvalidUInt8 = New("Result is not a valid uint8", 0, "")

// VMExecutionError parity returns this when there was an error executing the call in the VM
var VMExecutionError = New("VM execution error", 0, "")

// New returns a new rpcError
func New(err string, code int, details string) error {
	return &RpcError{
		err:     err,
		Code:    code,
		Details: details,
	}
}
