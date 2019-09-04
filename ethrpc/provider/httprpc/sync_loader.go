package httprpc

import "github.com/alethio/web3-go/jsonrpc2"

// SyncLoader is a synchronous loader that makes one http request per RPC
type SyncLoader struct {
	// this method provides the data for the loader
	fetch func(keys *jsonrpc2.JSONRPCRequest) ([]byte, error)
}

// NewSyncLoader creates a new syncLoader given a fetch, wait, and maxBatch
func NewSyncLoader() (*SyncLoader, error) {
	return &SyncLoader{}, nil
}

// Init initializes the SyncLoader
func (l *SyncLoader) Init(p *HTTPProvider) {
	l.fetch = p.fetchSingle
}

// Load turns a RPCRequest into a byte array response
func (l *SyncLoader) Load(req *jsonrpc2.JSONRPCRequest) ([]byte, error) {
	return l.fetch(req)
}
