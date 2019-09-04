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

// RawBalances : the tree-like structure representing all un-parsed balances
type RawBalances map[BlockNumber]map[Address]map[Source]string

// IntBalances : the tree-like structure with all balances converted to big.Int
type IntBalances map[BlockNumber]map[Address]map[Source]*big.Int

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

// BalanceResponse : the response associated with a balance request
type BalanceResponse struct {
	Request *BalanceRequest
	Balance string
}
