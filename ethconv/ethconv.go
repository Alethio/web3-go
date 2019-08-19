package ethconv

import (
	"fmt"
	"math/big"
	"strings"
)

// multipliers for various eth denominations
const (
	Wei  = "1"
	Gwei = "1000000000"
	Eth  = "1000000000000000000"
)

// ERC20 function signatures
const (
	ERC20Transfer     = "a9059cbb"
	ERC20TransferFrom = "23b872dd"
	/*
		0x06fdde03 -> [ function ] name
		0x095ea7b3 -> [ function ] approve
		0x18160ddd -> [ function ] totalSupply
		0x313ce567 -> [ function ] decimals
		0x475a9fa9 -> [ function ] issueTokens
		0x70a08231 -> [ function ] balanceOf
		0x95d89b41 -> [ function ] symbol
		0xdd62ed3e -> [ function ] allowance
		0xddf252ad -> [ event ] Transfer
		0x8c5be1e5 -> [ event ] Approval
	*/
)

// FromWei returns a string converted from wei to respective unit
func FromWei(value string, to string, prec uint) (string, bool) {
	v, ok := new(big.Float).SetString(value)
	if !ok {
		return "", false
	}
	unit, ok := new(big.Float).SetString(to)
	if !ok {
		return "", false
	}
	v.Quo(v, unit)

	return fmt.Sprintf("%.*f", prec, v), true
}

// IsERC20Transfer return true if the payload starts with erc20 transfer/transferFrom
func IsERC20Transfer(payload string) bool {
	if strings.HasPrefix(payload, ERC20Transfer) {
		return true
	}
	if strings.HasPrefix(payload, ERC20TransferFrom) {
		return true
	}
	return false
}
