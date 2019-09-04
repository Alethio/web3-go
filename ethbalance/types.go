package ethbalance

import (
	"math/big"

	"github.com/alethio/web3-go/ethrpc"
)

// Bookkeeper wraps the operations
type Bookkeeper struct {
	eth     ethrpc.ETHInterface
	retries uint
}

// BlockNumber : wrapper type for a block numnber
type BlockNumber uint64

// Address : wrapper type for an ETH address
type Address string

// RawBalanceSheet : the tree-like structure representing all un-parsed balances
type RawBalanceSheet map[BlockNumber]map[Address]map[Source]string

// IntBalanceSheet : the tree-like structure with all balances converted to big.Int
type IntBalanceSheet map[BlockNumber]map[Address]map[Source]*big.Int

// Source : either "ETH" or a token address
type Source string

const (
	// ETH : Ethereum source
	ETH Source = "ETH"
)

// BalanceRequest : a unit of work
type BalanceRequest struct {
	Block   BlockNumber
	Address Address
	Source  Source
}

// RawBalanceResponse : the raw response associated with a balance request
type RawBalanceResponse struct {
	Request *BalanceRequest
	Balance string
}

// IntBalanceResponse : the converted response associated with a balance request
type IntBalanceResponse struct {
	Request *BalanceRequest
	Balance *big.Int
}
