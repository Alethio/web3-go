package provider

import "encoding/json"

// Interface represents a web3 connection provider interface
type Interface interface {
	Start() error
	Stop()
	Call(result interface{}, method string, params ...interface{}) error
	CallRaw(method string, params ...interface{}) ([]byte, error)
	Subscribe(receiver chan *json.RawMessage, method string, event string, params ...interface{}) error
}
