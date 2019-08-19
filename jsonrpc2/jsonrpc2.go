package jsonrpc2

import "encoding/json"

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeClientRequest(method string, args interface{}, id string) ([]byte, error) {
	c := &JSONRPCRequest{
		Version: "2.0",
		Method:  method,
		Params:  args,
		ID:      id,
	}
	return json.Marshal(c)
}

// DecodeResponse decodes the top level json rpc response
func DecodeResponse(response []byte) (*JSONRPCMessage, error) {
	var msg JSONRPCMessage
	msg.Raw = response
	if err := json.Unmarshal(response, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
