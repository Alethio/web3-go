package ethrpc

// json rpc methods
const (
	// parity
	ParitySubscribe           = "parity_subscribe"
	ParityPendingTransactions = "parity_pendingTransactions"

	// geth
	GETHTxPoolContent = "txpool_content"

	// net
	NetPeerCount = "net_peerCount"

	// web3
	WEB3ClientVersion = "web3_clientVersion"

	// eth
	ETHBlockNumber                      = "eth_blockNumber"
	ETHCall                             = "eth_call"
	ETHGetBalance                       = "eth_getBalance"
	ETHGetBlockByNumber                 = "eth_getBlockByNumber"
	ETHGetBlockTransactionCountByNumber = "eth_getBlockTransactionCountByNumber"
	ETHGetCode                          = "eth_getCode"
	ETHGetFilterChanges                 = "eth_getFilterChanges"
	ETHGetTransactionByHash             = "eth_getTransactionByHash"
	ETHGetTransactionReceipt            = "eth_getTransactionReceipt"
	ETHGetUncleByBlockHashAndIndex      = "eth_getUncleByBlockHashAndIndex"
	ETHGetUncleByBlockNumberAndIndex    = "eth_getUncleByBlockNumberAndIndex"
	ETHPendingTransactionFilter         = "eth_newPendingTransactionFilter"
	ETHSubscribe                        = "eth_subscribe"

	// trace
	TraceBlock                   = "trace_block"
	TraceReplayBlockTransactions = "trace_replayBlockTransactions"

	// eth pubsub
	ETHNewHeads               = "newHeads"
	ETHNewPendingTransactions = "newPendingTransactions"

	// consts
	ClientGETH   = "geth"
	ClientParity = "parity"
)

// ERC20 signatures
const (
	// functions
	NameFunction         = "0x06fdde03"
	ApproveFunction      = "0x095ea7b3" // mandatory
	TotalSupplyFunction  = "0x18160ddd" // mandatory
	TransferFromFunction = "0x23b872dd" // mandatory
	DecimalsFunction     = "0x313ce567"
	IssueTokensFunction  = "0x475a9fa9"
	BalanceOfFunction    = "0x70a08231" // mandatory
	SymbolFunction       = "0x95d89b41"
	TransferFunction     = "0xa9059cbb" // mandatory
	AllowanceFunction    = "0xdd62ed3e" // mandatory

	// events
	TransferEvent = "0xddf252ad" // mandatory
	ApprovalEvent = "0x8c5be1e5" // mandatory
)

const (
	// DefaultCallGas is the default gas to use for eth_calls
	DefaultCallGas = "0xffffff"
)
