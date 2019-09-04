package httprpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alethio/web3-go/etherr"
	"github.com/alethio/web3-go/jsonrpc2"
)

const (
	// DefaultHTTPTimeout is the default timeout interval for http requests
	DefaultHTTPTimeout = 3 * time.Second
)

// HTTPProvider implements ethereum RPC calls over HTTP
type HTTPProvider struct {
	client      *http.Client
	url         string
	loader      RPCLoader
	httpTimeout time.Duration
}

type RPCLoader interface {
	Load(*jsonrpc2.JSONRPCRequest) ([]byte, error)
}

// Start does nothing on the http provider
func (p *HTTPProvider) Start() error {
	// TODO: maybe check if server is reachable?
	return nil
}

// Stop does nothing on the http provider
func (p *HTTPProvider) Stop() {
	return
}

// CallRaw calls a RPC method and returns the raw result
func (p *HTTPProvider) CallRaw(method string, params ...interface{}) ([]byte, error) {
	req := jsonrpc2.BuildRequest(method, params)
	return p.loader.Load(req)
}

// Call calls a RPC method and returns coresponding object
func (p *HTTPProvider) Call(result interface{}, method string, params ...interface{}) error {
	req := jsonrpc2.BuildRequest(method, params)
	raw, err := p.loader.Load(req)
	if err != nil {
		return fmt.Errorf("call: %s", err)
	}

	resp, err := jsonrpc2.DecodeResponse(raw)
	if err != nil {
		return fmt.Errorf("decode rpc message: %s", err)
	}

	null := string(json.RawMessage([]byte("null")))
	if string(resp.Result) == null {
		return etherr.Nil
	}

	if resp.Error != nil {
		switch resp.Error.Code {
		case -32015: // VM execution error
			err := etherr.VMExecutionError.(*etherr.RpcError)
			err.Code = resp.Error.Code
			err.Details = resp.Error.Data
			return err
		default:
			return etherr.New(resp.Error.Message, resp.Error.Code, resp.Error.Data)
		}
	}

	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return err
	}
	return nil

}

// Subscribe creates a subscription to event using method. not available on http
func (p *HTTPProvider) Subscribe(receiver chan *json.RawMessage, method string, event string, params ...interface{}) error {
	return fmt.Errorf("subscriptions not supported over http, please use websockets")
}

// New initializes a Client and returns it
func New(url string) (*HTTPProvider, error) {
	var httpClient = &http.Client{Transport: &http.Transport{}}

	p := &HTTPProvider{
		url:         url,
		client:      httpClient,
		httpTimeout: DefaultHTTPTimeout,
	}

	p.loader = newSyncLoader(syncLoaderConfig{
		Fetch: p.fetchSingle,
	})

	return p, nil
}

// NewWithBatch initializez a Client with batching and returns it
func NewWithBatch(url string, batchMaxSize int, batchWait time.Duration) (*HTTPProvider, error) {
	if batchMaxSize < 0 {
		return nil, fmt.Errorf("Maximum batch size can not be negative")
	}
	if batchWait < 1*time.Millisecond {
		return nil, fmt.Errorf("Minimum wait time must be at least 1 Millisecond")

	}

	var httpClient = &http.Client{Transport: &http.Transport{}}

	p := &HTTPProvider{
		url:         url,
		client:      httpClient,
		httpTimeout: DefaultHTTPTimeout,
	}

	p.loader = newBatchLoader(batchLoaderConfig{
		Wait:     batchWait,
		MaxBatch: batchMaxSize,
		Fetch:    p.fetchMultiple,
	})

	return p, nil
}

// SetHTTPTimeout allows setting the http timeout from outside
func (p *HTTPProvider) SetHTTPTimeout(httpTimeout time.Duration) {
	p.httpTimeout = httpTimeout
}
