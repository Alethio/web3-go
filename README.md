## Examples

```
go run main.go --eth-client-url wss://mainnet.infura.io/ws getBlockNumber
go run main.go --eth-client-url ws://alethio-geth-trace:9546 getCode 0xb8c77482e45f1f44de1745f52c74426c631bdd52
```


## structs
This package aims at defining a single source of data structures matching some of the parity JSONRPC responses

It currently supports the following calls: 
- `eth_getBlockByNumber`
- `eth_getTransactionReceipt`
- `eth_getUncleByBlockHashAndIndex`
- `trace_block`
- `trace_replayBlockTransactions`

## validator
This tool is intended for validating the logical integrity of JSONRPC responses coming from parity.

Checking one block would involve the following steps:
1. instantiate a new validator instance
2. load all the JSONRPC responses into the validator
3. call the `Run()` function which returns a boolean and an error

For more details, check the [example function](/validator/validator_test.go)
