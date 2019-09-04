package httprpc

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/alethio/web3-go/jsonrpc2"
	"github.com/sirupsen/logrus"
)

func (p *HTTPProvider) fetchSingle(request *jsonrpc2.JSONRPCRequest) ([]byte, error) {
	payload, err := request.Encode()
	if err != nil {
		return nil, err
	}

	return p.fetch(payload)
}

func (p *HTTPProvider) fetchMultiple(requests []*jsonrpc2.JSONRPCRequest) ([][]byte, []error) {
	payload, err := jsonrpc2.EncodeClientRequests(requests)
	if err != nil {
		return nil, []error{err}
	}

	response, err := p.fetch(payload)
	if err != nil {
		return [][]byte{}, []error{err}
	}

	responses, err := jsonrpc2.DecodeResponses(response)
	castedResponses := make([][]byte, len(responses))
	for index, resp := range responses {
		castedResponses[index] = []byte(resp)
	}
	return castedResponses, []error{err}
}

func (p *HTTPProvider) fetch(payload []byte) ([]byte, error) {
	logrus.Debugf("Making request with payload: %s\n", payload)

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
