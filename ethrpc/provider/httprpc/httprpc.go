package httprpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/alethio/web3-go/etherr"
	"github.com/alethio/web3-go/jsonrpc2"
)

type HTTPProvider struct {
	url    string
	client *http.Client
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
	return p.makeRequest(method, params)
}

// Call calls a RPC method and returns coresponding object
func (p *HTTPProvider) Call(result interface{}, method string, params ...interface{}) error {
	response, err := p.makeRequest(method, params)
	if err != nil {
		return fmt.Errorf("call: %s", err)
	}

	resp, err := jsonrpc2.DecodeResponse(response)
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
		url:    url,
		client: httpClient,
	}

	return p, nil
}

func (p *HTTPProvider) makeRequest(method string, params []interface{}) ([]byte, error) {
	id := strconv.FormatInt(rand.Int63(), 16)
	rpcRequest, err := jsonrpc2.EncodeClientRequest(method, params, id)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("POST", p.url, bytes.NewReader(rpcRequest))
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
