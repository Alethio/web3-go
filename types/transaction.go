package types

type Transaction struct {
	BlockHash        string      `json:"blockHash"`
	BlockNumber      string      `json:"blockNumber"`
	ChainId          string      `json:"chainId"`
	Condition        interface{} `json:"condition"`
	Creates          string      `json:"creates"`
	From             string      `json:"from"`
	Gas              string      `json:"gas"`
	GasPrice         string      `json:"gasPrice"`
	Hash             string      `json:"hash"`
	Input            string      `json:"input"`
	Nonce            string      `json:"nonce"`
	PublicKey        string      `json:"publicKey"`
	R                string      `json:"r"`
	Raw              string      `json:"raw"`
	S                string      `json:"s"`
	StandardV        string      `json:"standardV"`
	To               string      `json:"to"`
	TransactionIndex string      `json:"transactionIndex"`
	V                string      `json:"v"`
	Value            string      `json:"value"`
}
