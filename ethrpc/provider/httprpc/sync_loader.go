package httprpc

import "github.com/alethio/web3-go/jsonrpc2"

// syncLoaderConfig captures the config to create a new syncLoader
type syncLoaderConfig struct {
	// Fetch is a method that provides the data for the loader
	Fetch func(req *jsonrpc2.JSONRPCRequest) ([]byte, error)
}

type syncLoader struct {
	// this method provides the data for the loader
	fetch func(keys *jsonrpc2.JSONRPCRequest) ([]byte, error)
}

// newSyncLoader creates a new syncLoader given a fetch, wait, and maxBatch
func newSyncLoader(config syncLoaderConfig) *syncLoader {
	return &syncLoader{
		fetch: config.Fetch,
	}
}

// Load a request, batching will be applied automatically
func (l *syncLoader) Load(req *jsonrpc2.JSONRPCRequest) ([]byte, error) {
	return l.fetch(req)
}
