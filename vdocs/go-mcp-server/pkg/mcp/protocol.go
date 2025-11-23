package mcp

import (
	"encoding/json"
	"fmt"
)

// CreateRequest creates a JSON-RPC request message
func CreateRequest(id interface{}, method string, params interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

// CreateResponse creates a JSON-RPC response message
func CreateResponse(id interface{}, result interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// CreateErrorResponse creates a JSON-RPC error response
func CreateErrorResponse(id interface{}, err *RPCError) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}
}

// CreateNotification creates a JSON-RPC notification
func CreateNotification(method string, params interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
}

// ParseMessage parses a JSON-RPC message from bytes
func ParseMessage(data []byte) (*JSONRPCMessage, error) {
	var msg JSONRPCMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}
	
	// Validate JSON-RPC version
	if msg.JSONRPC != "2.0" {
		return nil, fmt.Errorf("invalid JSON-RPC version: %s", msg.JSONRPC)
	}
	
	return &msg, nil
}

// EncodeMessage encodes a JSON-RPC message to bytes
func EncodeMessage(msg *JSONRPCMessage) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}
	return data, nil
}

// IsRequest checks if a message is a request
func (m *JSONRPCMessage) IsRequest() bool {
	return m.Method != "" && m.ID != nil
}

// IsNotification checks if a message is a notification
func (m *JSONRPCMessage) IsNotification() bool {
	return m.Method != "" && m.ID == nil
}

// IsResponse checks if a message is a response
func (m *JSONRPCMessage) IsResponse() bool {
	return m.Method == "" && m.ID != nil
}

// HasError checks if a message has an error
func (m *JSONRPCMessage) HasError() bool {
	return m.Error != nil
}

