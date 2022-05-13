package pithy

import "fmt"

// APIError holds the error code and error message for responding to request peer
type APIResult struct {
	Message    string `json:"message"`
	ResultCode int    `json:"code"`
	// For internal use
	StatusCode int    `json:"-"`
	RawBytes   []byte `json:"-"` // If raw is not nil, directly send raw bytes rather than sending json bytes of the result
}

func (r *APIResult) String() string {
	return fmt.Sprintf("Result:%d Status:%d Message:%s", r.ResultCode, r.StatusCode, r.Message)
}

func NewAPIResult(result, status int, msg string) *APIResult {
	return &APIResult{
		ResultCode: result,
		StatusCode: status,
		Message:    msg,
	}
}

func NewAPIResultFromResultCode(result int, msg string) *APIResult {
	return NewAPIResult(result, 0, msg)
}
