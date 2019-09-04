package jsonrpc2

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeClientRequest(method string, args interface{}, id string) ([]byte, error) {
	return NewRequest(method, args, id).Encode()
}

// Encode returns the byte array json encoded value
func (req *JSONRPCRequest) Encode() ([]byte, error) {
	return json.Marshal(req)

}

// EncodeClientRequests encodes an array of rpc requests
func EncodeClientRequests(requests []*JSONRPCRequest) ([]byte, error) {
	return json.Marshal(requests)
}

// BuildRequest creates a new RPC request struct with a random ID
func BuildRequest(method string, args interface{}) *JSONRPCRequest {
	id := strconv.FormatInt(rand.Int63(), 16)
	return NewRequest(method, args, id)
}

// NewRequest creates a new RPC requests struct with all attributes required
func NewRequest(method string, args interface{}, id string) *JSONRPCRequest {
	return &JSONRPCRequest{
		Version: "2.0",
		Method:  method,
		Params:  args,
		ID:      id,
	}
}

// DecodeResponse decodes the top level json rpc response for a single rpc call
func DecodeResponse(response []byte) (*JSONRPCMessage, error) {
	var message JSONRPCMessage
	message.Raw = response
	if err := json.Unmarshal(response, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

// DecodeResponses decodes the top level json rpc response for a batched rpc call
func DecodeResponses(response []byte) ([]json.RawMessage, error) {
	var messages []json.RawMessage
	if err := json.Unmarshal(response, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
