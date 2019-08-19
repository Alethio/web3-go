package ethbalance

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/alethio/web3-go/ethrpc"

	"github.com/go-test/deep"
)

type MockETH struct {
	ethrpc.ETH
	balances   map[string]map[string]map[string]*big.Int
	throwError bool
}

func (m *MockETH) GetBalanceAtBlock(address, block string) (*big.Int, error) {
	if m.throwError == true {
		return nil, fmt.Errorf("Fatal error")
	}
	return m.balances[block][address]["eth"], nil
}

func (m *MockETH) GetTokenBalanceAtBlock(address, token, block string) (*big.Int, error) {
	if m.throwError == true {
		return nil, fmt.Errorf("Fatal error")
	}
	return m.balances[block][address][token], nil
}

func ExampleBookkeeper_GetBalancesAtBlock() {
	r, err := ethrpc.NewWithDefaults("wss://mainnet.infura.io/ws")
	if err != nil {
		fmt.Println(err)
		return
	}

	b := New(r, 10)
	accounts := map[string][]string{
		"0xa838e871a02c6d883bf004352fc7dac8f781fed6": []string{
			"0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1",
			"0x0f5d2fb29fb7d3cfee444a200298f468908cc942",
			"0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
			"0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2",
		},
	}
	balances, err := b.GetBalancesAtBlock(accounts, "0x62d313")

	if err != nil {
		fmt.Println(err)
		return
	}
	for account, accountBalances := range balances {
		for source, value := range accountBalances {
			fmt.Printf("%s[%s]: %v\n", account, source, value)
		}
	}
	// Example output:
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0xd26114cd6EE289AccF82350c8d8487fedB8A0C07]: 409757565152676909
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2]: 3600000000000000000000
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0x0f5d2fb29fb7d3cfee444a200298f468908cc942]: 7041922408306145321820
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[eth]: 1000670436501076869
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1]: 33780620000000000000
}

func TestGetBalancesWithOneAddressAndNoTokens(t *testing.T) {
	mockEth := &MockETH{}
	mockEth.balances = make(map[string]map[string]map[string]*big.Int)
	mockEth.balances["latest"] = make(map[string]map[string]*big.Int)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"] = make(map[string]*big.Int)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"]["eth"] = big.NewInt(100)

	bookkeeper := New(mockEth, 10)
	query := make(map[string][]string)
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"] = make([]string, 0)

	balances, err := bookkeeper.GetBalancesAtBlock(query, "latest")
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if diff := deep.Equal(balances, mockEth.balances["latest"]); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithOneAddressAndTokens(t *testing.T) {
	mockEth := &MockETH{}
	mockEth.balances = make(map[string]map[string]map[string]*big.Int)
	mockEth.balances["latest"] = make(map[string]map[string]*big.Int)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"] = make(map[string]*big.Int)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"]["eth"] = big.NewInt(100)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"]["0xabc"] = big.NewInt(102)

	bookkeeper := New(mockEth, 10)
	query := make(map[string][]string)
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"] = []string{"0xabc"}

	balances, err := bookkeeper.GetBalancesAtBlock(query, "latest")
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if diff := deep.Equal(balances, mockEth.balances["latest"]); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithMultipleAddressesAndTokens(t *testing.T) {
	mockEth := &MockETH{}
	mockEth.balances = make(map[string]map[string]map[string]*big.Int)
	mockEth.balances["latest"] = make(map[string]map[string]*big.Int)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"] = make(map[string]*big.Int)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"]["eth"] = big.NewInt(100)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"]["0xabc"] = big.NewInt(102)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"]["0xabd"] = big.NewInt(105)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af52"] = make(map[string]*big.Int)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af52"]["eth"] = big.NewInt(101)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af52"]["0xabc"] = big.NewInt(103)
	mockEth.balances["latest"]["0x9fc201b6bc40cccbd5b588532ce98b845f95af52"]["0xabf"] = big.NewInt(104)

	bookkeeper := New(mockEth, 10)
	query := make(map[string][]string)
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"] = []string{"0xabc", "0xabd"}
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af52"] = []string{"0xabc", "0xabf"}

	balances, err := bookkeeper.GetBalancesAtBlock(query, "latest")
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if diff := deep.Equal(balances, mockEth.balances["latest"]); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithError(t *testing.T) {
	mockEth := &MockETH{throwError: true}
	bookkeeper := New(mockEth, 10)
	query := make(map[string][]string)
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af51"] = []string{"0xabc", "0xabd"}
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af52"] = []string{"0xabc", "0xabf"}
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af53"] = []string{"0xabc", "0xabf"}
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af54"] = []string{"0xabc", "0xabf"}
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af55"] = []string{"0xabc", "0xabf"}
	query["0x9fc201b6bc40cccbd5b588532ce98b845f95af52"] = []string{"0xabc", "0xabf"}

	_, err := bookkeeper.GetBalancesAtBlock(query, "latest")
	if err == nil {
		t.Fatal("Expecting error")
	}
}
