package jsonrpc2

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// JSONRPCError : standard JSONRPC Error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// JSONRPCRequest standard JSONRPC 2 Request
type JSONRPCRequest struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      string      `json:"id"`
}

// JSONRPCMessage : standard JSONRPC 2 Response
type JSONRPCMessage struct {
	Version string          `json:"jsonrpc"`
	ID      string          `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Raw     json.RawMessage `json:"-"`
}

func (msg *JSONRPCMessage) IsNotification() bool {
	return msg.ID == "" && msg.Method != ""
}

func (msg *JSONRPCMessage) IsResponse() bool {
	return msg.HasValidID() && msg.Method == "" && len(msg.Params) == 0
}

func (msg *JSONRPCMessage) HasValidID() bool {
	return len(msg.ID) > 0 && msg.ID[0] != '{' && msg.ID[0] != '['
}

// ValidID decods the id if it s valid
func (msg *JSONRPCMessage) ValidID() (string, error) {
	if !msg.HasValidID() {
		return "", fmt.Errorf("Message does not have a valid id")
	}
	return string(msg.ID), nil
}

// UINTResult decods the the result as uint
func (msg *JSONRPCMessage) UINTResult() (uint64, error) {
	var s string
	err := json.Unmarshal(msg.Result, &s)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(s, 0, 64)
}

func (msg *JSONRPCMessage) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

// JSONRPCNotification : standard JSONRPC 2 Notification
type JSONRPCNotification struct {
	ID     string          `json:"subscription"`
	Result json.RawMessage `json:"result"`
}

// ValidID decods the id if it s valid
func (msg *JSONRPCNotification) ValidID() (string, error) {
	if len(msg.ID) == 0 {
		return "", fmt.Errorf("No ID found")
	}
	return msg.ID, nil
}
