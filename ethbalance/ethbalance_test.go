package ethbalance

import (
	"fmt"
	"testing"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/strhelper"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
)

type MockETH struct {
	ethrpc.ETH
	balances   RawBalanceSheet
	throwError bool
}

func (m *MockETH) GetRawBalanceAtBlock(address, block string) (string, error) {
	if m.throwError == true {
		return "", fmt.Errorf("Fatal error")
	}
	blockNumber, _ := strhelper.HexStrToInt64(block)
	return m.balances[BlockNumber(blockNumber)][address][ETH], nil
}

func (m *MockETH) GetRawTokenBalanceAtBlock(address, token, block string) (string, error) {
	if m.throwError == true {
		return "", fmt.Errorf("Fatal error")
	}
	blockNumber, _ := strhelper.HexStrToInt64(block)
	return m.balances[BlockNumber(blockNumber)][address][Source(token)], nil
}

func ExampleBookkeeper_GetIntBalanceResults() {
	r, err := ethrpc.NewWithDefaults("wss://mainnet.infura.io/ws")
	if err != nil {
		fmt.Println(err)
		return
	}

	b := New(r, 10)
	results, err := b.GetIntBalanceResults(balanceRequests())

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, res := range results {
		fmt.Printf("%s[%s]: %s\n", res.Request.Address, res.Request.Source, res.Balance)
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
	block := BlockNumber(7500000)
	mockEth.balances = RawBalanceSheet{
		block: map[Address]map[Source]string{
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af51": map[Source]string{
				ETH: fmt.Sprintf("0x%x", 100),
			},
		},
	}

	bookkeeper := New(mockEth, 10)
	requests := []*BalanceRequest{
		&BalanceRequest{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Source:  ETH,
			Block:   block,
		},
	}

	balances, err := bookkeeper.GetRawBalanceSheet(requests)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	spew.Dump(mockEth.balances)
	spew.Dump(balances)
	if diff := deep.Equal(balances, mockEth.balances); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithOneAddressAndTokens(t *testing.T) {
	mockEth := &MockETH{}
	block := BlockNumber(7500000)
	mockEth.balances = make(RawBalanceSheet)
	mockEth.balances = RawBalanceSheet{
		block: map[Address]map[Source]string{
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af51": map[Source]string{
				ETH:     fmt.Sprintf("0x%x", 100),
				"0xabc": fmt.Sprintf("0x%x", 102),
			},
		},
	}

	bookkeeper := New(mockEth, 10)
	requests := []*BalanceRequest{
		&BalanceRequest{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Source:  ETH,
			Block:   block,
		},
		&BalanceRequest{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Source:  "0xabc",
			Block:   block,
		},
	}

	balances, err := bookkeeper.GetRawBalanceSheet(requests)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if diff := deep.Equal(balances, mockEth.balances); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithMultipleAddressesAndTokens(t *testing.T) {
	mockEth := &MockETH{}
	block := BlockNumber(7500000)
	mockEth.balances = make(RawBalanceSheet)
	mockEth.balances = RawBalanceSheet{
		block: map[Address]map[Source]string{
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af51": map[Source]string{
				ETH:     fmt.Sprintf("0x%x", 100),
				"0xabc": fmt.Sprintf("0x%x", 102),
				"0xabd": fmt.Sprintf("0x%x", 105),
			},
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af52": map[Source]string{
				ETH:     fmt.Sprintf("0x%x", 101),
				"0xabc": fmt.Sprintf("0x%x", 103),
				"0xabd": fmt.Sprintf("0x%x", 104),
			},
		},
	}

	bookkeeper := New(mockEth, 10)
	requests := make([]*BalanceRequest, 0, 0)
	for _, address := range []Address{"0x9fc201b6bc40cccbd5b588532ce98b845f95af51", "0x9fc201b6bc40cccbd5b588532ce98b845f95af52"} {
		for _, source := range []Source{ETH, "0xabc", "0xabd"} {
			requests = append(requests, &BalanceRequest{
				Address: address,
				Source:  source,
				Block:   block,
			})
		}
	}

	balances, err := bookkeeper.GetRawBalanceSheet(requests)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if diff := deep.Equal(balances, mockEth.balances); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithError(t *testing.T) {
	mockEth := &MockETH{throwError: true}
	bookkeeper := New(mockEth, 10)
	block := BlockNumber(7500000)
	requests := []*BalanceRequest{
		&BalanceRequest{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Source:  ETH,
			Block:   block,
		},
		&BalanceRequest{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Source:  "0xabc",
			Block:   block,
		},
	}

	_, err := bookkeeper.GetRawBalanceResults(requests)
	if err == nil {
		t.Fatal("Expecting error")
	}
}

func balanceRequests() []*BalanceRequest {
	address := Address("0xa838e871a02c6d883bf004352fc7dac8f781fed6")
	block := BlockNumber(7500000)
	return []*BalanceRequest{
		&BalanceRequest{
			Address: address,
			Block:   block,
			Source:  ETH,
		},
		&BalanceRequest{
			Address: address,
			Block:   block,
			Source:  Source("0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1"),
		},
		&BalanceRequest{
			Address: address,
			Block:   block,
			Source:  Source("0x0f5d2fb29fb7d3cfee444a200298f468908cc942"),
		},
		&BalanceRequest{
			Address: address,
			Block:   block,
			Source:  Source("0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"),
		},
		&BalanceRequest{
			Address: address,
			Block:   block,
			Source:  Source("0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2"),
		},
	}
}
