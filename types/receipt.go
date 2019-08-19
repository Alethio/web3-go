package types

type Receipt struct {
	BlockHash         string      `json:"blockHash"`
	BlockNumber       string      `json:"blockNumber"`
	ContractAddress   interface{} `json:"contractAddress"`
	CumulativeGasUsed string      `json:"cumulativeGasUsed"`
	From              string      `json:"from"`
	GasUsed           string      `json:"gasUsed"`
	Logs              []Log       `json:"logs"`
	LogsBloom         string      `json:"logsBloom"`
	Root              string      `json:"root"`
	Status            string      `json:"status"`
	To                string      `json:"to"`
	TransactionHash   string      `json:"transactionHash"`
	TransactionIndex  string      `json:"transactionIndex"`
}

type Log struct {
	Address             string   `json:"address"`
	BlockHash           string   `json:"blockHash"`
	BlockNumber         string   `json:"blockNumber"`
	Data                string   `json:"data"`
	LogIndex            string   `json:"logIndex"`
	Removed             bool     `json:"removed"`
	Topics              []string `json:"topics"`
	TransactionHash     string   `json:"transactionHash"`
	TransactionIndex    string   `json:"transactionIndex"`
	TransactionLogIndex string   `json:"transactionLogIndex"`
	Type                string   `json:"type"`
}

type RPCGetTransactionReceipt struct {
	Jsonrpc string  `json:"jsonrpc"`
	Result  Receipt `json:"result"`
	ID      int     `json:"id"`
}
