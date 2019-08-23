package validator

import (
	"github.com/alethio/web3-go/types"
)

// Validator is intended for validating the logical integrity of JSONRPC responses coming from parity
type Validator struct {
	Block types.Block
	Uncles []types.Block
	Receipts []types.Receipt
	Traces []types.Trace
	Replays []types.TransactionReplay
	
	loadedMap map[string]bool
}

// New returns a new Validator instance
func New() *Validator {
	return &Validator{loadedMap:make(map[string]bool)}
}

// Run executes all the available verifiers and returns (true, nil) if the block is valid
// or (false, error) if the block is not valid
func (v *Validator) Run() (bool, error) {
	err := v.verifyBlock()
	if err != nil {
		return false, err
	}

	err = v.verifyUncles()
	if err != nil {
		return false, err
	}

	err = v.verifyReceipts()
	if err != nil {
		return false, err
	}

	err = v.verifyTrace()
	if err != nil {
		return false, err
	}

	err = v.verifyReplay()
	if err != nil {
		return false, err
	}

	return true, nil
}
