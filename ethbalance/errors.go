package ethbalance

import "fmt"

// RequestError wraps a request with the error received
type RequestError struct {
	Request *BalanceRequest
	Err     error
}

func (re RequestError) section() string {
	return fmt.Sprintf("    %d %s %s", re.Request.Block, re.Request.Address, re.Request.Source)
}

// CollectBalancesError wraps errors returned from the RPC requests
type CollectBalancesError struct {
	Errors []*RequestError
}

// Error message aggregates all unique errors
func (cbe CollectBalancesError) Error() string {
	fullErrorMessage := renderFullErrorMessage(cbe.Errors)
	return fmt.Sprintf("Unable to collect balances because of these errors:\n%s", fullErrorMessage)
}

// DecodeBalancesError wraps errors returned when trying to decode the responses
type DecodeBalancesError struct {
	Errors []*RequestError
}

// Error message aggregates all unique errors
func (dbe DecodeBalancesError) Error() string {
	fullErrorMessage := renderFullErrorMessage(dbe.Errors)
	return fmt.Sprintf("Unable to collect balances because of these errors:\n%s", fullErrorMessage)
}

func renderFullErrorMessage(errors []*RequestError) string {
	var fullErrorMessage string
	errorMessages := make(map[string]string)
	index := 0

	for _, reqError := range errors {
		errorMessageKey := reqError.Err.Error()
		sectionError, ok := errorMessages[errorMessageKey]
		if !ok {
			sectionError = fmt.Sprintf("[%d] %s\n    Requests:\n", index, errorMessageKey)
			index++
		}
		sectionError += reqError.section() + "\n"
		errorMessages[errorMessageKey] = sectionError
	}

	for _, errorMessage := range errorMessages {
		fullErrorMessage += errorMessage
	}

	return fullErrorMessage
}
