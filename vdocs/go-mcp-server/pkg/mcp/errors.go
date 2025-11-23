package mcp

import "fmt"

// NewRPCError creates a new RPC error
func NewRPCError(code int, message string, data interface{}) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewParseError creates a parse error
func NewParseError(message string) *RPCError {
	return NewRPCError(ParseErrorCode, "Parse error: "+message, nil)
}

// NewInvalidRequestError creates an invalid request error
func NewInvalidRequestError(message string) *RPCError {
	return NewRPCError(InvalidRequestCode, "Invalid request: "+message, nil)
}

// NewMethodNotFoundError creates a method not found error
func NewMethodNotFoundError(method string) *RPCError {
	return NewRPCError(MethodNotFoundCode, fmt.Sprintf("Method not found: %s", method), nil)
}

// NewInvalidParamsError creates an invalid params error
func NewInvalidParamsError(message string) *RPCError {
	return NewRPCError(InvalidParamsCode, "Invalid params: "+message, nil)
}

// NewInternalError creates an internal error
func NewInternalError(message string) *RPCError {
	return NewRPCError(InternalErrorCode, "Internal error: "+message, nil)
}

// Error implements the error interface
func (e *RPCError) Error() string {
	return fmt.Sprintf("RPC Error %d: %s", e.Code, e.Message)
}

