// Inspired by https://github.com/vektah/dataloaden with some tweaks
package httprpc

import (
	"sync"
	"time"

	"github.com/alethio/web3-go/jsonrpc2"
)

// batchLoaderConfig captures the config to create a new batchLoader
type batchLoaderConfig struct {
	// Fetch is a method that provides the data for the loader
	Fetch func(req []*jsonrpc2.JSONRPCRequest) ([][]byte, []error)

	// Wait is how long wait before sending a batch
	Wait time.Duration

	// MaxBatch will limit the maximum number of keys to send in one batch, 0 = not limit
	MaxBatch int
}

// newBatchLoader creates a new batchLoader given a fetch, wait, and maxBatch
func newBatchLoader(config batchLoaderConfig) *batchLoader {
	return &batchLoader{
		fetch:    config.Fetch,
		wait:     config.Wait,
		maxBatch: config.MaxBatch,
	}
}

// batchLoader batches and caches requests
type batchLoader struct {
	// this method provides the data for the loader
	fetch func(keys []*jsonrpc2.JSONRPCRequest) ([][]byte, []error)

	// how long to done before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// INTERNAL

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *batchLoaderBatch

	// mutex to prevent races
	mu sync.Mutex
}

type batchLoaderBatch struct {
	requests []*jsonrpc2.JSONRPCRequest
	data     [][]byte
	error    []error
	closing  bool
	done     chan struct{}
}

// Load a request, batching will be applied automatically
func (l *batchLoader) Load(req *jsonrpc2.JSONRPCRequest) ([]byte, error) {
	return l.LoadThunk(req)()
}

// LoadThunk returns a function that when called will block waiting for a byte.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *batchLoader) LoadThunk(req *jsonrpc2.JSONRPCRequest) func() ([]byte, error) {
	l.mu.Lock()
	if l.batch == nil {
		l.batch = &batchLoaderBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.reqIndex(l, req)
	l.mu.Unlock()

	return func() ([]byte, error) {
		<-batch.done

		var data []byte
		if pos < len(batch.data) {
			data = batch.data[pos]
		}

		var err error
		// its convenient to be able to return a single error for everything like timeout
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		return data, err
	}
}

// LoadAll fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
func (l *batchLoader) LoadAll(reqs []*jsonrpc2.JSONRPCRequest) ([][]byte, []error) {
	results := make([]func() ([]byte, error), len(reqs))

	for i, req := range reqs {
		results[i] = l.LoadThunk(req)
	}

	bytes := make([][]byte, len(reqs))
	errors := make([]error, len(reqs))
	for i, thunk := range results {
		bytes[i], errors[i] = thunk()
	}
	return bytes, errors
}

// LoadAllThunk returns a function that when called will block waiting for a bytes.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *batchLoader) LoadAllThunk(keys []*jsonrpc2.JSONRPCRequest) func() ([][]byte, []error) {
	results := make([]func() ([]byte, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]byte, []error) {
		bytes := make([][]byte, len(keys))
		errors := make([]error, len(keys))
		for i, thunk := range results {
			bytes[i], errors[i] = thunk()
		}
		return bytes, errors
	}
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *batchLoaderBatch) reqIndex(l *batchLoader, req *jsonrpc2.JSONRPCRequest) int {
	for i, existingRequest := range b.requests {
		if req == existingRequest {
			return i
		}
	}

	pos := len(b.requests)
	b.requests = append(b.requests, req)
	if pos == 0 {
		go b.startTimer(l)
	}

	if l.maxBatch != 0 && pos >= l.maxBatch-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}

	return pos
}

func (b *batchLoaderBatch) startTimer(l *batchLoader) {
	time.Sleep(l.wait)
	l.mu.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if b.closing {
		l.mu.Unlock()
		return
	}

	l.batch = nil
	l.mu.Unlock()

	b.end(l)
}

func (b *batchLoaderBatch) end(l *batchLoader) {
	b.data, b.error = l.fetch(b.requests)
	close(b.done)
}
