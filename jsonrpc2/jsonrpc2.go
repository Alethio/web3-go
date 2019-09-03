package jsonrpc2

import "encoding/json"

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeClientRequest(method string, args interface{}, id string) ([]byte, error) {
	c := NewJSONRPCRequest(method, args, id)
	return json.Marshal(c)
}

// EncodeClientRequests encodes an array of rpc requests
func EncodeClientRequests(requests []*JSONRPCRequest) ([]byte, error) {
	return json.Marshal(requests)
}

// NewJSONRPCRequest creates a new RPC request struct
func NewJSONRPCRequest(method string, args interface{}, id string) *JSONRPCRequest {
	return &JSONRPCRequest{
		Version: "2.0",
		Method:  method,
		Params:  args,
		ID:      id,
	}
}

// DecodeResponse decodes the top level json rpc response
func DecodeResponse(response []byte) (*JSONRPCMessage, error) {
	var msg *JSONRPCMessage
	if err := json.Unmarshal(response, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// DecodeResponses decodes the top level json rpc response
func DecodeResponses(response []byte) ([]*JSONRPCMessage, error) {
	var msgs []*JSONRPCMessage
	if err := json.Unmarshal(response, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
