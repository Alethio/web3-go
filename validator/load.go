package validator

import "encoding/json"

func (v *Validator) LoadBlockResponse(data []byte) error {
	return json.Unmarshal(data, &v.ResponseBlock)
}

func (v *Validator) LoadUnclesResponse(data []byte) error {
	return json.Unmarshal(data, &v.ResponseUncles)
}

func (v *Validator) LoadReceiptsResponse(data []byte) error {
	return json.Unmarshal(data, &v.ResponseReceipts)
}

func (v *Validator) LoadTraceBlockResponse(data []byte) error {
	return json.Unmarshal(data, &v.ResponseTrace)
}

func (v *Validator) LoadReplayResponse(data []byte) error {
	return json.Unmarshal(data, &v.ResponseReplay)
}
