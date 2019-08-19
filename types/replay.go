package types

type TransactionReplay struct {
	Output          string                 `json:"output"`
	StateDiff       map[string]interface{} `json:"stateDiff"`
	Trace           []Trace                `json:"trace"`
	VMTrace         *VMTrace               `json:"vmTrace"`
	TransactionHash *string                `json:"transactionHash"`
}

type VMTrace struct {
	Code string      `json:"code"`
	Ops  []VMTraceOp `json:"ops"`
}

type VMTraceOp struct {
	Cost int `json:"cost"`
	Ex   struct {
		Mem   interface{} `json:"mem"`
		Push  interface{} `json:"push"`
		Store interface{} `json:"store"`
		Used  int         `json:"used"`
	} `json:"ex"`
	Pc  int      `json:"pc"`
	Sub *VMTrace `json:"sub"`
}

type RPCTraceReplayBlockTransactions struct {
	Jsonrpc string              `json:"jsonrpc"`
	Result  []TransactionReplay `json:"result"`
	ID      int                 `json:"id"`
}
