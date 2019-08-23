package validator

import "github.com/alethio/web3-go/types"

func (v *Validator) LoadBlock(block types.Block) {
	v.Block = block
	v.loadedMap[Block] = true
}

func (v *Validator) LoadUncles(uncles []types.Block) {
	v.Uncles = uncles
	v.loadedMap[Uncles] = true
}

func (v *Validator) LoadReceipts(receipts []types.Receipt) {
	v.Receipts = receipts
	v.loadedMap[Receipts] = true
}

func (v *Validator) LoadTraces(traces []types.Trace) {
	v.Traces = traces
	v.loadedMap[Traces] = true
}

func (v *Validator) LoadReplays(replays []types.TransactionReplay) {
	v.Replays = replays
	v.loadedMap[Replays] = true
}
