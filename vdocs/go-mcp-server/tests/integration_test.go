package tests

import (
	"encoding/json"
	"testing"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/server"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
)

// TestServerCreation tests that server can be created successfully
func TestServerCreation(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AccessKeyID = "test-key"
	cfg.AccessKeySecret = "test-secret"
	
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	if srv == nil {
		t.Fatal("Expected server to be non-nil")
	}
}

// TestHandlerInitialize tests the initialize handler
func TestHandlerInitialize(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AccessKeyID = "test-key"
	cfg.AccessKeySecret = "test-secret"
	
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	handler := server.NewHandler(srv)
	
	// Create initialize request
	initReq := mcp.InitializeRequest{
		ProtocolVersion: mcp.ProtocolVersion,
		ClientInfo: mcp.ClientInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
	}
	
	paramsData, _ := json.Marshal(initReq)
	var params interface{}
	json.Unmarshal(paramsData, &params)
	
	msg := mcp.CreateRequest(1, "initialize", params)
	
	// Handle the request
	response, err := handler.Handle(msg)
	if err != nil {
		t.Fatalf("Handler error: %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.ID != msg.ID {
		t.Errorf("Expected response ID %v, got %v", msg.ID, response.ID)
	}
	
	if response.Result == nil {
		t.Fatal("Expected result in response")
	}
}

// TestHandlerToolsList tests the tools/list handler
func TestHandlerToolsList(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AccessKeyID = "test-key"
	cfg.AccessKeySecret = "test-secret"
	
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	// Register tools
	if err := srv.GetRegistry().RegisterAllTools(); err != nil {
		t.Fatalf("Failed to register tools: %v", err)
	}
	
	handler := server.NewHandler(srv)
	
	msg := mcp.CreateRequest(2, "tools/list", nil)
	
	response, err := handler.Handle(msg)
	if err != nil {
		t.Fatalf("Handler error: %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.Result == nil {
		t.Fatal("Expected result in response")
	}
	
	// Parse result
	resultData, _ := json.Marshal(response.Result)
	var listResult mcp.ListToolsResult
	if err := json.Unmarshal(resultData, &listResult); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}
	
	if len(listResult.Tools) == 0 {
		t.Error("Expected at least one tool to be registered")
	}
	
	t.Logf("Registered tools count: %d", len(listResult.Tools))
}

// TestHandlerPing tests the ping handler
func TestHandlerPing(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AccessKeyID = "test-key"
	cfg.AccessKeySecret = "test-secret"
	
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	handler := server.NewHandler(srv)
	
	msg := mcp.CreateRequest(3, "ping", nil)
	
	response, err := handler.Handle(msg)
	if err != nil {
		t.Fatalf("Handler error: %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.Result == nil {
		t.Fatal("Expected result in response")
	}
}

// TestHandlerMethodNotFound tests handling of unknown methods
func TestHandlerMethodNotFound(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AccessKeyID = "test-key"
	cfg.AccessKeySecret = "test-secret"
	
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	handler := server.NewHandler(srv)
	
	msg := mcp.CreateRequest(4, "unknown_method", nil)
	
	response, err := handler.Handle(msg)
	if err != nil {
		t.Fatalf("Handler error: %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.Error == nil {
		t.Fatal("Expected error in response")
	}
	
	if response.Error.Code != mcp.MethodNotFoundCode {
		t.Errorf("Expected error code %d, got %d", mcp.MethodNotFoundCode, response.Error.Code)
	}
}

