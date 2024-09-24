package model

// CreateAPILogParams holds the parameters for logging API actions related to user creation.
type CreateAPILogParams struct {
	Method       string
	RequestData  interface{}
	ResponseData interface{}
}
