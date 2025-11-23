package tests

import (
	"log"
	"os"
	"testing"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/tools"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
)

func TestRegistryCreation(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := log.New(os.Stderr, "[Test] ", log.LstdFlags)
	
	registry := tools.NewRegistry(cfg, logger)
	if registry == nil {
		t.Fatal("Expected registry to be created")
	}
}

func TestRegisterTool(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := log.New(os.Stderr, "[Test] ", log.LstdFlags)
	registry := tools.NewRegistry(cfg, logger)
	
	tool := mcp.Tool{
		Name:        "TestTool",
		Description: "A test tool",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
	}
	
	toolFunc := func(args map[string]interface{}) (interface{}, error) {
		return "test result", nil
	}
	
	registry.RegisterTool(tool, toolFunc)
	
	toolsList := registry.ListTools()
	if len(toolsList) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(toolsList))
	}
	
	if toolsList[0].Name != "TestTool" {
		t.Errorf("Expected tool name TestTool, got %s", toolsList[0].Name)
	}
}

func TestExecuteTool(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := log.New(os.Stderr, "[Test] ", log.LstdFlags)
	registry := tools.NewRegistry(cfg, logger)
	
	tool := mcp.Tool{
		Name:        "TestTool",
		Description: "A test tool",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
	}
	
	toolFunc := func(args map[string]interface{}) (interface{}, error) {
		return map[string]string{"status": "success"}, nil
	}
	
	registry.RegisterTool(tool, toolFunc)
	
	result, err := registry.ExecuteTool("TestTool", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to execute tool: %v", err)
	}
	
	if result.IsError {
		t.Error("Expected successful execution")
	}
	
	if len(result.Content) == 0 {
		t.Error("Expected non-empty content")
	}
}

func TestExecuteNonExistentTool(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := log.New(os.Stderr, "[Test] ", log.LstdFlags)
	registry := tools.NewRegistry(cfg, logger)
	
	_, err := registry.ExecuteTool("NonExistent", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error for non-existent tool")
	}
}

