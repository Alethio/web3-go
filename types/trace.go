package types

type TraceAction struct {
	// call
	CallType *string `json:"callType"`
	To       *string `json:"to"`
	Input    *string `json:"input"`
	// call + create
	From  *string `json:"from"`
	Gas   *string `json:"gas"`
	Value *string `json:"value"`
	// create
	Init *string `json:"init"`
	// suicide
	Address       *string `json:"address"`
	Balance       *string `json:"balance"`
	RefundAddress *string `json:"refundAddress"`
}

type TraceResult struct {
	// call
	Output *string `json:"output"`
	// call + create
	GasUsed *string `json:"gasUsed"`
	// create
	Address *string `json:"address"`
	Code    *string `json:"code"`
	// suicide is nil
}

type Trace struct {
	Action              TraceAction  `json:"action"`
	BlockHash           *string      `json:"blockHash"`
	BlockNumber         *int         `json:"blockNumber"`
	Result              *TraceResult `json:"result"`
	Subtraces           int          `json:"subtraces"`
	TraceAddress        []int        `json:"traceAddress"`
	TransactionHash     *string      `json:"transactionHash"`
	TransactionPosition *int         `json:"transactionPosition"`
	Type                string       `json:"type"`
	Error               *string      `json:"error"`
}

type RPCTraceBlock struct {
	Jsonrpc string  `json:"jsonrpc"`
	Result  []Trace `json:"result"`
	ID      int     `json:"id"`
}
