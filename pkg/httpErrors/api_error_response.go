package httpErrors

// ApiErrorResponse acts as a struct dto for error response
type ApiErrorResponse struct {
	ErrorInfo ErrorInfo `json:"error"`
	RequestId string    `json:"request_id"`
}

type ErrorInfo struct {
	Status   int                   `json:"status"`
	Message  string                `json:"message"`
	SubError *[]ApiValidationError `json:"sub_errors,omitempty"`
}

type ApiValidationError struct {
	Field         string `json:"field"`
	RejectedValue any    `json:"rejected_value"`
	Message       string `json:"message"`
}

// NewApiErrorResponse is a factory function
func NewApiErrorResponse(errInfo ErrorInfo, requestId string) ApiErrorResponse {
	return ApiErrorResponse{
		ErrorInfo: errInfo,
		RequestId: requestId,
	}
}
