package tests

import (
	"testing"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
)

func TestCreateRequest(t *testing.T) {
	req := mcp.CreateRequest(1, "initialize", map[string]string{"test": "value"})
	
	if req.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", req.JSONRPC)
	}
	
	if req.ID != 1 {
		t.Errorf("Expected ID 1, got %v", req.ID)
	}
	
	if req.Method != "initialize" {
		t.Errorf("Expected method initialize, got %s", req.Method)
	}
}

func TestCreateResponse(t *testing.T) {
	resp := mcp.CreateResponse(1, map[string]string{"result": "success"})
	
	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", resp.JSONRPC)
	}
	
	if resp.ID != 1 {
		t.Errorf("Expected ID 1, got %v", resp.ID)
	}
	
	if resp.Result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestCreateErrorResponse(t *testing.T) {
	err := mcp.NewInternalError("test error")
	resp := mcp.CreateErrorResponse(1, err)
	
	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", resp.JSONRPC)
	}
	
	if resp.Error == nil {
		t.Error("Expected error to be non-nil")
	}
	
	if resp.Error.Code != mcp.InternalErrorCode {
		t.Errorf("Expected error code %d, got %d", mcp.InternalErrorCode, resp.Error.Code)
	}
}

func TestParseMessage(t *testing.T) {
	jsonData := []byte(`{"jsonrpc":"2.0","id":1,"method":"test","params":{}}`)
	
	msg, err := mcp.ParseMessage(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse message: %v", err)
	}
	
	if msg.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", msg.JSONRPC)
	}
	
	if msg.Method != "test" {
		t.Errorf("Expected method test, got %s", msg.Method)
	}
}

func TestEncodeMessage(t *testing.T) {
	msg := mcp.CreateRequest(1, "test", nil)
	
	data, err := mcp.EncodeMessage(msg)
	if err != nil {
		t.Fatalf("Failed to encode message: %v", err)
	}
	
	if len(data) == 0 {
		t.Error("Expected non-empty encoded data")
	}
}

func TestIsRequest(t *testing.T) {
	req := mcp.CreateRequest(1, "test", nil)
	if !req.IsRequest() {
		t.Error("Expected message to be identified as request")
	}
	
	resp := mcp.CreateResponse(1, nil)
	if resp.IsRequest() {
		t.Error("Expected message not to be identified as request")
	}
}

func TestIsNotification(t *testing.T) {
	notif := mcp.CreateNotification("test", nil)
	if !notif.IsNotification() {
		t.Error("Expected message to be identified as notification")
	}
	
	req := mcp.CreateRequest(1, "test", nil)
	if req.IsNotification() {
		t.Error("Expected message not to be identified as notification")
	}
}

func TestIsResponse(t *testing.T) {
	resp := mcp.CreateResponse(1, nil)
	if !resp.IsResponse() {
		t.Error("Expected message to be identified as response")
	}
	
	req := mcp.CreateRequest(1, "test", nil)
	if req.IsResponse() {
		t.Error("Expected message not to be identified as response")
	}
}

