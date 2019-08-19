package types

type BlockHeader struct {
	Author           string        `json:"author"`
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	SealFields       []string      `json:"sealFields"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TransactionsRoot string        `json:"transactionsRoot"`
}

type Block struct {
	BlockHeader
	Size             string        `json:"size"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []Transaction `json:"transactions"`
	Uncles           []string      `json:"uncles"`
}


// TODO check this
type RPCGetBlockByNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Block  `json:"result"`
	ID      int    `json:"id"`
}

// TODO check this
type RPCGetUncleByBlockHashAndIndex struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Block  `json:"result"`
	ID      int    `json:"id"`
}
