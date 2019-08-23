package validator

import (
	"encoding/json"

	"github.com/alethio/web3-go/types"
)

func (v *Validator) LoadBlockResponse(data []byte) error {
	var r types.RPCGetBlockByNumberResponse

	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	v.Block = r.Result
	v.loadedMap[Block] = true

	return nil
}

func (v *Validator) LoadUnclesResponse(data []byte) error {
	var r []types.RPCGetUncleByBlockHashAndIndex

	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	for _, u := range r {
		v.Uncles = append(v.Uncles, u.Result)
	}
	v.loadedMap[Uncles] = true

	return nil
}

func (v *Validator) LoadReceiptsResponse(data []byte) error {
	var r []types.RPCGetTransactionReceipt

	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	for _, rec := range r {
		v.Receipts = append(v.Receipts, rec.Result)
	}

	v.loadedMap[Receipts] = true

	return nil
}

func (v *Validator) LoadTraceBlockResponse(data []byte) error {
	var r types.RPCTraceBlock

	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	v.Traces = r.Result
	v.loadedMap[Traces] = true

	return nil
}

func (v *Validator) LoadReplayResponse(data []byte) error {
	var r types.RPCTraceReplayBlockTransactions

	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	v.Replays = r.Result
	v.loadedMap[Replays] = true

	return nil
}
