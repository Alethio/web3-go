package validator

import (
	"fmt"
	"strconv"
	"strings"
)

func (v *Validator) isLoaded(item string) bool {
	_, exists := v.loadedMap[item]
	return exists
}

func (v *Validator) verifyBlock() error {
	if !v.isLoaded(Block) {
		return fmt.Errorf("block is mandatory")
	}

	if v.Block.Hash == "" {
		return fmt.Errorf("block hash is empty")
	}

	return nil
}

func (v *Validator) verifyUncles() error {
	if !v.isLoaded(Uncles) {
		return nil
	}

	if len(v.Uncles) != len(v.Block.Uncles) {
		return fmt.Errorf("uncles count is different")
	}

	for i, hash := range v.Block.Uncles {
		if v.Uncles[i].Hash != hash {
			return fmt.Errorf("uncle hash at index %d does not match", i)
		}
	}

	return nil
}

func (v *Validator) verifyReceipts() error {
	if !v.isLoaded(Receipts) {
		return nil
	}

	if len(v.Receipts) != len(v.Block.Transactions) {
		return fmt.Errorf("receipts count is different")
	}

	for i, receipt := range v.Receipts {
		tx := v.Block.Transactions[i]
		r := receipt

		if r.TransactionHash != tx.Hash {
			return fmt.Errorf("receipt at index %d does not match transaction hash", i)
		}

		if r.TransactionIndex != tx.TransactionIndex {
			return fmt.Errorf("receipt at index %d does not match transaction index", i)
		}

		if r.BlockHash != tx.BlockHash {
			return fmt.Errorf("receipt at index %d does not match block hash", i)
		}

		if r.BlockNumber != tx.BlockNumber {
			return fmt.Errorf("receipt at index %d does not match block number", i)
		}
	}

	return nil
}

func (v *Validator) verifyTrace() error {
	if !v.isLoaded(Traces) {
		return nil
	}

	if !v.isLoaded(Receipts) {
		return fmt.Errorf("receipts are mandatory to validate traces")
	}

	blockHash := v.Block.Hash
	blockNumberInt64, err := strconv.ParseInt(strings.TrimPrefix(v.Block.Number, "0x"), 16, 64)
	if err != nil {
		return err
	}

	blockNumber := int(blockNumberInt64)
	uniqueTransactions := make(map[int]string, len(v.Receipts))

	for i, trace := range v.Traces {
		if trace.Type != "reward" {
			uniqueTransactions[*trace.TransactionPosition] = *trace.TransactionHash

			if v.Block.Transactions[*trace.TransactionPosition].Hash != *trace.TransactionHash {
				return fmt.Errorf("trace at index %d does not match transaction hash at position %d", i, *trace.TransactionPosition)
			}
		}

		if trace.BlockNumber != nil && *trace.BlockNumber != blockNumber {
			return fmt.Errorf("trace at index %d does not match block number", i)
		}

		if trace.BlockHash != nil && *trace.BlockHash != blockHash {
			return fmt.Errorf("trace at index %d does not match block hash", i)
		}
	}

	// verify that each transaction present on the block has at least one coresponding trace
	for k, tx := range v.Block.Transactions {
		if uniqueTransactions[k] != tx.Hash {
			return fmt.Errorf("did not find any trace for transaction %s", tx.Hash)
		}
	}

	return nil
}

func (v *Validator) verifyReplay() error {
	if !v.isLoaded(Replays) {
		return nil
	}

	if len(v.Replays) != len(v.Block.Transactions) {
		return fmt.Errorf("replay count does not match transaction count")
	}

	for i, replay := range v.Replays {
		tx := v.Block.Transactions[i]

		if replay.TransactionHash != nil {
			if *replay.TransactionHash != tx.Hash {
				return fmt.Errorf("replay at index %d does not match transaction hash", i)
			}
		}

		if len(replay.Trace) == 0 {
			return fmt.Errorf("replay at index %d has empty trace", i)
		}

		// The first trace in the replay.Trace should represent the transaction itself
		firstTrace := replay.Trace[0]

		if firstTrace.Action.From == nil || *firstTrace.Action.From != tx.From {
			return fmt.Errorf("replay at index %d field 'from' does not match transaction", i)
		}

		// fixme: don't check this because it looks like the gas is totally skewed in the traces/replays.
		// if firstTrace.Action.Gas == nil || *firstTrace.Action.Gas != tx.Gas {
		// 	return fmt.Errorf("replay at index %d field 'gas' does not match transaction", i)
		// }

		if firstTrace.Action.Value == nil || *firstTrace.Action.Value != tx.Value {
			return fmt.Errorf("replay at index %d field 'value' does not match transaction", i)
		}

		switch firstTrace.Type {
		case "create":
			if firstTrace.Action.Init == nil ||
				*firstTrace.Action.Init != tx.Input {
				return fmt.Errorf("replay at index %d for a 'create' does not match input of transaction", i)
			}
		case "call":
			if firstTrace.Action.Input == nil ||
				*firstTrace.Action.Input != tx.Input {
				return fmt.Errorf("replay at index %d for a 'call' does not match input of transaction", i)
			}

			if firstTrace.Action.To == nil ||
				*firstTrace.Action.To != tx.To {
				return fmt.Errorf("replay at index %d field 'to' does not match transaction", i)
			}
		default:
			return fmt.Errorf("replay at index %d: invalid transaction type", i)
		}
	}

	return nil
}
