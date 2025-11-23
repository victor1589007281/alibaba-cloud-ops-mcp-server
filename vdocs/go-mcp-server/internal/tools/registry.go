package tools

import (
	"encoding/json"
	"fmt"
	"log"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/alibabacloud"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
)

// ToolFunc is a function that executes a tool
type ToolFunc func(args map[string]interface{}) (interface{}, error)

// ToolDefinition represents a tool definition
type ToolDefinition struct {
	Tool mcp.Tool
	Func ToolFunc
}

// Registry manages tool registration and execution
type Registry struct {
	tools       map[string]*ToolDefinition
	cfg         *config.Config
	credsProv   *alibabacloud.CredentialsProvider
	apiMeta     *alibabacloud.APIMetaClient
	logger      *log.Logger
}

// NewRegistry creates a new tool registry
func NewRegistry(cfg *config.Config, logger *log.Logger) *Registry {
	return &Registry{
		tools:     make(map[string]*ToolDefinition),
		cfg:       cfg,
		credsProv: alibabacloud.NewCredentialsProvider(cfg),
		apiMeta:   alibabacloud.NewAPIMetaClient(),
		logger:    logger,
	}
}

// RegisterTool registers a tool
func (r *Registry) RegisterTool(tool mcp.Tool, fn ToolFunc) {
	r.tools[tool.Name] = &ToolDefinition{
		Tool: tool,
		Func: fn,
	}
	r.logger.Printf("Registered tool: %s", tool.Name)
}

// ListTools returns all registered tools
func (r *Registry) ListTools() []mcp.Tool {
	tools := make([]mcp.Tool, 0, len(r.tools))
	for _, def := range r.tools {
		tools = append(tools, def.Tool)
	}
	return tools
}

// ExecuteTool executes a tool by name
func (r *Registry) ExecuteTool(name string, args map[string]interface{}) (mcp.CallToolResult, error) {
	def, ok := r.tools[name]
	if !ok {
		return mcp.CallToolResult{}, fmt.Errorf("tool not found: %s", name)
	}
	
	r.logger.Printf("Executing tool: %s with args: %v", name, args)
	
	// Execute the tool function
	result, err := def.Func(args)
	if err != nil {
		return mcp.CallToolResult{
			Content: []mcp.Content{
				{
					Type: "text",
					Text: fmt.Sprintf("Error executing tool: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	
	// Convert result to JSON string
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.CallToolResult{
			Content: []mcp.Content{
				{
					Type: "text",
					Text: fmt.Sprintf("Error marshaling result: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	
	return mcp.CallToolResult{
		Content: []mcp.Content{
			{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		IsError: false,
	}, nil
}

// RegisterAllTools registers all available tools
func (r *Registry) RegisterAllTools() error {
	// Register OOS tools
	if err := RegisterOOSTools(r); err != nil {
		return fmt.Errorf("failed to register OOS tools: %w", err)
	}
	
	// Register CMS tools
	if err := RegisterCMSTools(r); err != nil {
		return fmt.Errorf("failed to register CMS tools: %w", err)
	}
	
	// Register OSS tools
	if err := RegisterOSSTools(r); err != nil {
		return fmt.Errorf("failed to register OSS tools: %w", err)
	}
	
	// Register API tools
	if err := RegisterAPITools(r); err != nil {
		return fmt.Errorf("failed to register API tools: %w", err)
	}
	
	return nil
}

// GetCredentials gets credentials for API calls
func (r *Registry) GetCredentials() (*alibabacloud.Credentials, error) {
	return r.credsProv.GetCredentialsSimple()
}

// GetAPIMetaClient returns the API meta client
func (r *Registry) GetAPIMetaClient() *alibabacloud.APIMetaClient {
	return r.apiMeta
}

// GetConfig returns the configuration
func (r *Registry) GetConfig() *config.Config {
	return r.cfg
}

