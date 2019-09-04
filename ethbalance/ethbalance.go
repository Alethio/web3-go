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

// GetIntBalances takes a list of balance requests and returns a tree like
// structure containing all big.Int balances
func (b *Bookkeeper) GetIntBalances(requests []*BalanceRequest) (IntBalances, error) {
	balances := make(IntBalances)
	rawBalances, err := b.GetRawBalances(requests)
	if err != nil {
		return balances, err
	}

	for block, addresses := range rawBalances {
		balances[block] = make(map[Address]map[Source]*big.Int)
		for address, sources := range addresses {
			balances[block][address] = make(map[Source]*big.Int)
			for source, rawBalance := range sources {
				balance, err := strhelper.HexStrToBigInt(rawBalance)
				if err != nil {
					return balances, err
				}
				balances[block][address][source] = balance
			}
		}
	}
	return balances, nil

}

// GetRawBalances takes a list of balance requests and returns a tree like
// structure containing all string balances
func (b *Bookkeeper) GetRawBalances(requests []*BalanceRequest) (RawBalances, error) {
	balances := make(RawBalances)
	results := make(chan *BalanceResponse)
	done := make(chan error, 1)

	go b.fetchRequests(requests, results, done)

	for {
		select {
		case result := <-results:
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
		case err := <-done:
			return balances, err
		}
	}
}

func (b *Bookkeeper) fetchRequests(requests []*BalanceRequest, results chan *BalanceResponse, done chan error) {
	var tries uint = 0
	wg := sync.WaitGroup{}

	for {
		failed := make(chan *RequestError, len(requests))
		errors := make(chan error, len(requests))
		for _, request := range requests {
			wg.Add(1)
			go func(req *BalanceRequest, results chan *BalanceResponse, failed chan *RequestError) {
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
					results <- &BalanceResponse{
						Request: req,
						Balance: balance,
					}
					errors <- err
				}
			}(request, results, failed)
		}

		wg.Wait()
		close(failed)

		requests := make([]*BalanceRequest, 0, len(requests))
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
