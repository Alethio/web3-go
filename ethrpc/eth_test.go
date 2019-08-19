package ethrpc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	mock "github.com/alethio/ethmock/server"
)

func TestRequests(t *testing.T) {
	eth, teardown := setup(t)
	defer teardown()

	var tests = map[string]func(t *testing.T){
		// GetBalanceAtBlock(address, blockNumber string) (*big.Int, error)
		// GetBlockByNumber(number string, full bool) (b entities.RPCBlockResponse, err error)
		// GetBlockNumber() (int64, error)
		"GetBlockNumber": func(t *testing.T) {
			expected := int64(7912466)
			actual, err := eth.GetBlockNumber()

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		},
		// GetBlockTransactionCountByNumber(number string) (count string, err error)
		// GetClient() (string, error)
		"GetClient": func(t *testing.T) {
			expected := "geth"
			actual, err := eth.GetClient()

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		},
		// GetCode(a string) ([]byte, error)
		"GetCode - Empty": func(t *testing.T) {
			expected := []byte("")
			actual, err := eth.GetCode("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		},
		"GetCode": func(t *testing.T) {
			expected, _ := hex.DecodeString("6060604052361561004a576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063c8fea2fb1461004e578063ff03ad56146100af575b5b5b005b341561005957600080fd5b6100ad600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506100e6565b005b6100e4600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610210565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610209578290508073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85846000604051602001526040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15156101eb57600080fd5b6102c65a03f115156101fc57600080fd5b50505060405180519050505b5b5b50505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156102c857803373ffffffffffffffffffffffffffffffffffffffff16311015156102c6578173ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f1935050505015156102c557600080fd5b5b5b5b5b50505600a165627a7a72305820fe8e241acfb8e01e0521fa37f6b871952cdb705636e58db8766ec68e79742e930029")
			actual, err := eth.GetCode("0xa86c21c635273770c921c9ab4f966611fbe683a3")

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		},

		// GetContractName(address string) (string, error)
		// GetContractSymbol(address string) (string, error)
		// GetContractTotalSupply(address string) (*big.Int, error)
		// GetERC20Decimals(address string) (uint8, error)
		// GetFilterChanges(id string) (t []interface{}, err error)
		// GetLatestBlock(full bool) (b entities.RPCBlockResponse, err error)
		// GetPeerCount() (peers int64, err error)
		"GetPeerCount": func(t *testing.T) {
			expected := int64(100)
			actual, err := eth.GetPeerCount()

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		},
		// GetPendingFilterChanges(id string) (t []string, err error)
		// GetPendingTransactions() ([]entities.RPCTransactionResponse, error)
		// GetTokenBalanceAtBlock(address, token, blockNumber string) (*big.Int, error)
		// GetTransactionByHash(hash string) (entities.RPCTransactionResponse, error)
		// GetTransactionReceipt(hash string) (r entities.RPCReceiptResponse, err error)
		// GetUncleByBlockNumberAndIndex(blockNumber string, index string) (b entities.RPCBlockResponse, err error)
		// GetVersion() (ver string, err error)
		"GetVersion": func(t *testing.T) {
			expected := "Geth/v1.8.22-omnibus-260f7fbd/linux-amd64/go1.11.1"
			actual, err := eth.GetVersion()

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		},
	}

	for n, fn := range tests {
		t.Run(n, fn)
	}
}

func setup(t *testing.T) (*ETH, func() error) {
	t.Helper()

	srv, err := mock.New(8545, "../testdata/mock")
	assert.Nil(t, err)
	go srv.Serve()

	p, err := httprpc.New("http://localhost:8545")
	if err != nil {
		t.Fatal(err)
	}

	e, err := New(p)
	if err != nil {
		t.Fatal(err)
	}

	return e, srv.Close
}
