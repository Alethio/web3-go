package multicall

import (
	"encoding/hex"
	"fmt"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestViewCall(t *testing.T) {
	vc := ViewCall{
		Key: "key",
		Target: "0x0",
		Method: "balanceOf(address, uint64)(int256)",
		Arguments: []interface{}{"0x1234", uint64(12)},
	}
	expectedArgTypes := []string{"address", "uint64"}
	expectedCallData := []byte{
		0x29, 0x5e, 0xaa, 0xdf, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x12, 0x34, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0xc}
	assert.Equal(t, expectedArgTypes, vc.argumentTypes())
	callData, err := vc.callData()
	assert.Nil(t, err)
	assert.Equal(t, expectedCallData, callData)
}

func ExampleViwCall() {
	eth, err := getETH("https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d")
	vcs := ViewCalls{
		{
			Key:       "key.4",
			Target:    "0x5eb3fa2dfecdde21c950813c665e9364fa609bd2",
			Method:    "getLastBlockHash()(bytes32)",
			Arguments: []interface{}{},
		},
	}
	mc, _ := New(eth, Config{
		Preset: "mainnet",
	})
	block := "latest"
	res, err := mc.Call(vcs, block)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	blockHash := res.Calls["key.4"].ReturnValues[0].([32]byte)
	fmt.Println(hex.EncodeToString(blockHash[:]))
	fmt.Println(err)

}

func getETH(url string) (ethrpc.ETHInterface, error) {
	provider, err := httprpc.New(url)
	if err != nil {
		return nil, err
	}
	provider.SetHTTPTimeout(5 * time.Second)
	return ethrpc.New(provider)
}

