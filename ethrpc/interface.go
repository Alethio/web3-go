package ethrpc

import (
	"encoding/json"
	"math/big"

	"github.com/alethio/web3-go/types"
)

// ETHInterface defines the packages interface
type ETHInterface interface {
	CallContractFunction(function string, address string, gas string) (string, error)
	CallContractFunctionBigInt(function string, address string) (*big.Int, error)
	CallContractFunctionInt64(function string, address string) (int64, error)
	GetBalanceAtBlock(address, blockNumber string) (*big.Int, error)
	GetBlockByNumber(number string) (b types.Block, err error)
	GetBlockNumber() (int64, error)
	GetBlockTransactionCountByNumber(number string) (count string, err error)
	GetClient() (string, error)
	GetCode(a string) ([]byte, error)
	GetContractName(address string) (string, error)
	GetContractSymbol(address string) (string, error)
	GetContractTotalSupply(address string) (*big.Int, error)
	GetERC20Decimals(address string) (uint8, error)
	GetFilterChanges(id string) (t []interface{}, err error)
	GetLatestBlock() (b types.Block, err error)
	GetPeerCount() (peers int64, err error)
	GetPendingFilterChanges(id string) (t []string, err error)
	GetPendingTransactions() ([]types.Transaction, error)
	GetTokenBalanceAtBlock(address, token, blockNumber string) (*big.Int, error)
	GetTransactionByHash(hash string) (types.Transaction, error)
	GetTransactionReceipt(hash string) (r types.Receipt, err error)
	GetUncleByBlockHashAndIndex(hash string, index string) (b types.Block, err error)
	GetUncleByBlockNumberAndIndex(blockNumber string, index string) (b types.Block, err error)
	GetVersion() (ver string, err error)
	TraceBlock(blockNumber string) ([]types.Trace, error)
	TraceReplayBlockTransactions(blockNumber string, traceTypes ...string) ([]types.TransactionReplay, error)
	MakeRequest(result interface{}, method string, params ...interface{}) error
	NewBlockNumberSubscription() (r chan *int64, err error)
	NewHeadsSubscription() (r chan *types.BlockHeader, err error)
	NewPendingTransactionsSubscription() (r chan *string, err error)
	SetPendingTransactionsFilter() (id string, err error)
	Start() error
	Subscribe(receiver chan *json.RawMessage, method string, event string, params ...interface{}) error
}
