package ethbalance

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/strhelper"
)

// New returns a new Bookkeeper struct
func New(eth ethrpc.ETHInterface, retries uint) *Bookkeeper {
	return &Bookkeeper{
		eth:     eth,
		retries: retries,
	}
}

// GetIntBalanceSheet takes a list of balance requests and returns a tree like
// structure containing all int balances
func (b *Bookkeeper) GetIntBalanceSheet(requests []*BalanceRequest) (IntBalanceSheet, error) {
	balances := make(IntBalanceSheet)
	intResponses, err := b.GetIntBalanceResults(requests)
	for _, result := range intResponses {
		block := result.Request.Block
		address := result.Request.Address
		source := result.Request.Source

		if balances[block] == nil {
			balances[block] = make(map[Address]map[Source]*big.Int)
		}

		if balances[block][address] == nil {
			balances[block][address] = make(map[Source]*big.Int)

		}

		balances[block][address][source] = result.Balance
	}
	return balances, err
}

// GetRawBalanceSheet takes a list of balance requests and returns a tree like
// structure containing all hex string balances
func (b *Bookkeeper) GetRawBalanceSheet(requests []*BalanceRequest) (RawBalanceSheet, error) {
	balances := make(RawBalanceSheet)
	rawResponses, err := b.GetRawBalanceResults(requests)
	for _, result := range rawResponses {
		block := result.Request.Block
		address := result.Request.Address
		source := result.Request.Source

		if balances[block] == nil {
			balances[block] = make(map[Address]map[Source]string)
		}

		if balances[block][address] == nil {
			balances[block][address] = make(map[Source]string)

		}

		balances[block][address][source] = result.Balance
	}
	return balances, err
}

// GetIntBalanceResults returns an array of *big.Int balance results for the provided requests
func (b *Bookkeeper) GetIntBalanceResults(requests []*BalanceRequest) ([]*IntBalanceResponse, error) {
	intResponses := make([]*IntBalanceResponse, 0, len(requests))
	failedRequests := make([]*RequestError, 0, len(requests))

	rawResponses, err := b.GetRawBalanceResults(requests)
	if err != nil {
		if collectError, ok := err.(CollectBalancesError); ok {
			failedRequests = append(failedRequests, collectError.Errors...)
		} else {
			return intResponses, err
		}
	}

	for _, rawResponse := range rawResponses {
		intBalance, err := strhelper.HexStrToBigInt(rawResponse.Balance)
		if err != nil {
			failedRequests = append(failedRequests, &RequestError{rawResponse.Request, err})
		} else {
			intResponses = append(intResponses, &IntBalanceResponse{
				Request: rawResponse.Request,
				Balance: intBalance,
			})
		}
	}

	if len(failedRequests) > 0 {
		return intResponses, CollectBalancesError{failedRequests}
	}
	return intResponses, nil
}

// GetRawBalanceResults returns an array of hex string balance results for the provided requests
func (b *Bookkeeper) GetRawBalanceResults(requests []*BalanceRequest) ([]*RawBalanceResponse, error) {
	results := make(chan *RawBalanceResponse)
	responses := make([]*RawBalanceResponse, 0, len(requests))

	done := make(chan error, 1)

	go b.fetchRequests(requests, results, done)

	for {
		select {
		case result := <-results:
			responses = append(responses, result)
		case err := <-done:
			return responses, err
		}
	}

}

func (b *Bookkeeper) fetchRequests(requests []*BalanceRequest, results chan *RawBalanceResponse, done chan error) {
	var tries uint = 0
	wg := sync.WaitGroup{}

	for {
		failed := make(chan *RequestError, len(requests))
		errors := make(chan error, len(requests))
		for _, request := range requests {
			wg.Add(1)
			go func(req *BalanceRequest, results chan *RawBalanceResponse, failed chan *RequestError) {
				defer wg.Done()
				var balance string
				var err error

				block := fmt.Sprintf("0x%x", req.Block)
				address := string(req.Address)
				if req.Source == ETH {
					balance, err = b.eth.GetRawBalanceAtBlock(address, block)
				} else {
					token := string(req.Source)
					balance, err = b.eth.GetRawTokenBalanceAtBlock(address, token, block)
				}

				if err != nil {
					failed <- &RequestError{req, err}
				} else {
					results <- &RawBalanceResponse{
						Request: req,
						Balance: balance,
					}
					errors <- err
				}
			}(request, results, failed)
		}

		wg.Wait()
		close(failed)

		requests = make([]*BalanceRequest, 0, len(requests))
		reqErrors := make([]*RequestError, 0, len(requests))

		for reqError := range failed {
			reqErrors = append(reqErrors, reqError)
			requests = append(requests, reqError.Request)
		}

		if len(requests) == 0 {
			done <- nil
			return
		}

		tries++
		if tries >= b.retries {
			done <- CollectBalancesError{reqErrors}
			return
		}
	}
}
