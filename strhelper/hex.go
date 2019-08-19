package strhelper

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// HexStrToBigInt transforms a hex sting like "0xff" to a big.Int. Arbitrary length values are possible.
func HexStrToBigInt(hexString string) (*big.Int, error) {
	value := new(big.Int)
	_, ok := value.SetString(Trim0x(hexString), 16)
	if !ok {
		return value, fmt.Errorf("Could not transform hex string to big int: %s", hexString)
	}

	return value, nil
}

// HexStrToInt64 transforms a hex sting like "0xff" to an int64 like 15. The maximum value is "0x7fffffffffffffff", respectively 9223372036854775807
func HexStrToInt64(hexString string) (int64, error) {
	if !strings.HasPrefix(hexString, "0x") {
		hexString = "0x" + hexString
	}
	return strconv.ParseInt(hexString, 0, 64)
}
