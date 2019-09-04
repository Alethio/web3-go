package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/alethio/web3-go/ethbalance"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
)

type worker struct {
	eth ethrpc.ETHInterface
}

func main() {

	var ethURL string
	var batched bool
	flag.StringVar(&ethURL, "eth-client-url", "ws://localhost:8546", "Websockets URL of an Ethereum Client (parity needed)")
	flag.BoolVar(&batched, "batched", false, "Control wether the client is in batch mode")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Println("Please issue a command:")
		log.Println("  getBlockNumber")
		log.Println("  getCode address")
		log.Println("  newBlockNumberSubscription")
		os.Exit(0)
	}
	log.SetLevel(log.DebugLevel)

	var e *ethrpc.ETH
	var err error
	if batched {
		provider, err := httprpc.NewWithBatch(ethURL, 0, 4*time.Millisecond)
		if err != nil {
			log.Fatal(err)
		}
		e, err = ethrpc.New(provider)
	} else {
		e, err = ethrpc.NewWithDefaults(ethURL)
	}
	if err != nil {
		log.Fatal(err)
	}
	w := worker{
		eth: e,
	}

	cmd := args[0]
	switch cmd {
	case "getBlockNumber":
		n, err := w.eth.GetBlockNumber()
		if err != nil {
			log.Fatal("Eth failed to get block number: ", err)
		}

		log.Println(n)
	case "getLatestBlock":
		b, err := w.eth.GetLatestBlock()
		if err != nil {
			log.Fatal("Eth failed to get latest: ", err)
		}

		log.Printf("%+v\n", b)
	case "getBlockNumberRaw":
		ba, err := e.MakeRequestRaw(ethrpc.ETHBlockNumber)

		if err != nil {
			log.Fatal("Eth failed to get block number raw: ", err)
		}

		log.Println(string(ba))
	case "getCode":
		if len(args) < 2 {
			log.Fatal("Missing address")
		}
		ba, err := w.eth.GetCode(args[1])
		if err != nil {
			log.Fatal("Eth failed to get block number: ", err)
		}

		log.Println("code", ba)
	case "getUncleByBlockHashAndIndex":
		if len(args) < 3 {
			log.Fatal("Missing hash and or/index")
		}
		b, err := w.eth.GetUncleByBlockHashAndIndex(args[1], args[2])
		if err != nil {
			log.Fatal("Eth failed to get block number: ", err)
		}

		log.Printf("%+v\n", b)
	case "traceBlock":
		if len(args) < 2 {
			log.Fatal("Missing hex block number and/or trace types")
		}
		t, err := w.eth.TraceBlock(args[1])
		if err != nil {
			log.Fatal("Eth failed to trace block: ", err)
		}

		log.Printf("%+v\n", t)
	case "traceReplayBlockTransactions":
		if len(args) < 3 {
			log.Fatal("Missing hex block number")
		}
		r, err := w.eth.TraceReplayBlockTransactions(args[1], args[2:]...)
		if err != nil {
			log.Fatal("Eth failed to replay block: ", err)
		}

		log.Printf("%+v\n", r)
	case "newBlockNumberSubscription":
		blockNumbers, err := w.eth.NewBlockNumberSubscription()
		if err != nil {
			log.Fatal("Eth failed to get block number subscription: ", err)
		}

		// the subscription closes when the connection dies
		for number := range blockNumbers {
			log.Println(*number)
		}
		log.Warnf("subscription died")
	case "newHeadsSubscription":
		blockHeads, err := w.eth.NewHeadsSubscription()
		if err != nil {
			log.Fatal("Eth failed to get block number subscription: ", err)
		}

		// the subscription closes when the connection dies
		for head := range blockHeads {
			log.Printf("%+v\n", head)
		}
		log.Warnf("subscription died")
	case "getBalance":
		balance, err := w.eth.GetBalanceAtBlock(args[1], args[2])
		if err != nil {
			log.Fatal("Eth failed to get balance: ", err)
		}
		log.Println(balance)
	case "getTokenBalance":
		balance, err := w.eth.GetTokenBalanceAtBlock(args[1], args[2], args[3])
		if err != nil {
			log.Fatal("Eth failed to get balance: ", err)
		}
		log.Println(balance)
	case "getBalances":
		b := ethbalance.New(w.eth, 10)
		accounts := map[string][]string{
			"0xa838e871a02c6d883bf004352fc7dac8f781fed6": []string{
				"0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1",
				"0x0f5d2fb29fb7d3cfee444a200298f468908cc942",
				"0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
				"0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2",
			},
		}
		balances, err := b.GetBalancesAtBlock(accounts, "latest")
		if err != nil {
			fmt.Println(err)
			return
		}
		for account, accountBalances := range balances {
			for source, value := range accountBalances {
				fmt.Printf("%s[%s]: %v\n", account, source, value)
			}
		}
	default:
		log.Println("Command not implemented")
	}

}
