package httprpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
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
	batch       bool
	loader      *rpcLoader
	httpTimeout time.Duration
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
	return p.fetchRequestRaw(method, params)
}

// Call calls a RPC method and returns coresponding object
func (p *HTTPProvider) Call(result interface{}, method string, params ...interface{}) error {
	resp, err := p.makeRequest(method, params)
	if err != nil {
		return fmt.Errorf("call: %s", err)
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
		batch:       false,
		httpTimeout: DefaultHTTPTimeout,
	}

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
		batch:       true,
		httpTimeout: DefaultHTTPTimeout,
	}

	p.loader = newRPCLoader(rpcLoaderConfig{
		Wait:     batchWait,
		MaxBatch: batchMaxSize,
		Fetch:    p.fetchRequests,
	})

	return p, nil
}

// SetHTTPTimeout allows setting the http timeout from outside
func (p *HTTPProvider) SetHTTPTimeout(httpTimeout time.Duration) {
	p.httpTimeout = httpTimeout
}

func (p *HTTPProvider) makeRequest(method string, params []interface{}) (*jsonrpc2.JSONRPCMessage, error) {
	id := strconv.FormatInt(rand.Int63(), 16)
	req := jsonrpc2.NewJSONRPCRequest(method, params, id)

	if p.batch == true {
		return p.makeRequestAsync(req)
	} else {
		return p.makeRequestSync(req)
	}
}

func (p *HTTPProvider) makeRequestAsync(req *jsonrpc2.JSONRPCRequest) (*jsonrpc2.JSONRPCMessage, error) {
	return p.loader.Load(req)
}

func (p *HTTPProvider) makeRequestSync(req *jsonrpc2.JSONRPCRequest) (*jsonrpc2.JSONRPCMessage, error) {
	resp, err := p.fetchRequests([]*jsonrpc2.JSONRPCRequest{req})
	return resp[0], err[0]
}

func (p *HTTPProvider) fetchRequestRaw(method string, params []interface{}) ([]byte, error) {
	id := strconv.FormatInt(rand.Int63(), 16)
	payload, err := jsonrpc2.EncodeClientRequest(method, params, id)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("POST", p.url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer httpRequest.Body.Close()

	httpRequest.Header.Add("Content-Type", "application/json")

	response, err := p.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (p *HTTPProvider) fetchRequests(requests []*jsonrpc2.JSONRPCRequest) ([]*jsonrpc2.JSONRPCMessage, []error) {
	fmt.Println("Making http rpc request for:")
	fmt.Println(requests)

	payload, err := jsonrpc2.EncodeClientRequests(requests)
	if err != nil {
		return nil, []error{err}
	}

	httpRequest, err := http.NewRequest("POST", p.url, bytes.NewReader(payload))
	if err != nil {
		return nil, []error{err}
	}
	defer httpRequest.Body.Close()

	httpRequest.Header.Add("Content-Type", "application/json")

	response, err := p.client.Do(httpRequest)
	if err != nil {
		return nil, []error{err}
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, []error{err}
	}

	msgs, err := jsonrpc2.DecodeResponses(responseBody)
	if err != nil {
		return msgs, []error{fmt.Errorf("decode rpc message: %s", err)}
	}

	fmt.Println("Got back: ")
	fmt.Println(msgs)
	// TODO: ensure that the order of the msgs is the same as the requests in the payload
	return msgs, []error{nil}
}
