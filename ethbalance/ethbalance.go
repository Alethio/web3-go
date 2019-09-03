package ethbalance

import (
	"fmt"
	"math/big"

	"github.com/alethio/web3-go/ethrpc"

	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
)

// Bookkeeper provides highlevel access to ether and token balances.
type Bookkeeper struct {
	eth     ethrpc.ETHInterface
	workers int64
}

type balanceQuery struct {
	address      string
	forToken     bool
	tokenAddress string
	value        *big.Int
}

// New returns a new Bookkeeper struct
func New(eth ethrpc.ETHInterface, workers int) *Bookkeeper {
	return &Bookkeeper{
		eth:     eth,
		workers: int64(workers),
	}
}

// GetBalancesAtBlock gets a list of accounts and associated tokens and returns all balances
func (b *Bookkeeper) GetBalancesAtBlock(accounts map[string][]string, block string) (map[string]map[string]*big.Int, error) {
	balances := make(map[string]map[string]*big.Int)

	results := make(chan *balanceQuery)
	queries := make(chan *balanceQuery)
	errors := make(chan error)
	done := make(chan bool, 1)

	go b.scheduleQueries(accounts, queries)
	go b.processQueries(block, queries, results, errors, done)

	for {
		select {
		case result := <-results:
			if balances[result.address] == nil {
				balances[result.address] = make(map[string]*big.Int)
			}
			if result.forToken == true {
				balances[result.address][result.tokenAddress] = result.value
			} else {
				balances[result.address]["eth"] = result.value
			}
		case err := <-errors:
			return nil, err
		case <-done:
			return balances, nil
		}
	}
}

func (b *Bookkeeper) scheduleQueries(accounts map[string][]string, queries chan *balanceQuery) {
	for account, tokens := range accounts {
		queries <- &balanceQuery{
			address:      account,
			forToken:     false,
			tokenAddress: "",
			value:        nil,
		}

		for _, token := range tokens {
			queries <- &balanceQuery{
				address:      account,
				forToken:     true,
				tokenAddress: token,
				value:        nil,
			}
		}
	}
	close(queries)
}

func (b *Bookkeeper) processQueries(block string, queries, results chan *balanceQuery, errors chan error, done chan bool) {
	ctx := context.TODO()
	sem := semaphore.NewWeighted(b.workers)
	killSwitch := make(chan bool, 1)

	for {
		select {
		case <-killSwitch:
			return
		case query, ok := <-queries:
			if ok == false {
				if err := sem.Acquire(ctx, b.workers); err != nil {
					errors <- err
				} else {
					done <- true
				}
				return
			}

			if err := sem.Acquire(ctx, 1); err != nil {
				errors <- err
				killSwitch <- true
				return
			}
			go func(query *balanceQuery) {
				defer sem.Release(1)
				var balance *big.Int
				var err error

				if query.forToken == true {
					balance, err = b.eth.GetTokenBalanceAtBlock(query.address, query.tokenAddress, block)
				} else {
					balance, err = b.eth.GetBalanceAtBlock(query.address, block)
				}
				query.value = balance

				if err != nil {
					fmt.Println(err)
					killSwitch <- true
					errors <- err
				} else {
					results <- query
				}
			}(query)
		}
	}
}
