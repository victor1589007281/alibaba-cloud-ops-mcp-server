package server

import (
	"encoding/json"
	"fmt"
	"log"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
)

// Handler handles MCP protocol messages
type Handler struct {
	server *Server
	logger *log.Logger
}

// NewHandler creates a new handler
func NewHandler(server *Server) *Handler {
	return &Handler{
		server: server,
		logger: server.logger,
	}
}

// Handle handles an incoming message and returns a response
func (h *Handler) Handle(msg *mcp.JSONRPCMessage) (*mcp.JSONRPCMessage, error) {
	// Handle notifications (no response needed)
	if msg.IsNotification() {
		h.logger.Printf("Received notification: %s", msg.Method)
		return nil, nil
	}
	
	// Handle requests
	if msg.IsRequest() {
		return h.handleRequest(msg)
	}
	
	// Unknown message type
	return mcp.CreateErrorResponse(msg.ID, mcp.NewInvalidRequestError("invalid message type")), nil
}

// handleRequest handles a request message
func (h *Handler) handleRequest(msg *mcp.JSONRPCMessage) (*mcp.JSONRPCMessage, error) {
	switch msg.Method {
	case "initialize":
		return h.handleInitialize(msg)
	case "tools/list":
		return h.handleToolsList(msg)
	case "tools/call":
		return h.handleToolsCall(msg)
	case "ping":
		return h.handlePing(msg)
	default:
		return mcp.CreateErrorResponse(msg.ID, mcp.NewMethodNotFoundError(msg.Method)), nil
	}
}

// handleInitialize handles the initialize request
func (h *Handler) handleInitialize(msg *mcp.JSONRPCMessage) (*mcp.JSONRPCMessage, error) {
	h.logger.Println("Handling initialize request")
	
	result := mcp.InitializeResult{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities: mcp.Capabilities{
			Tools: &mcp.ToolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: mcp.ServerInfo{
			Name:    h.server.config.ServerName,
			Version: h.server.config.ServerVersion,
		},
	}
	
	return mcp.CreateResponse(msg.ID, result), nil
}

// handleToolsList handles the tools/list request
func (h *Handler) handleToolsList(msg *mcp.JSONRPCMessage) (*mcp.JSONRPCMessage, error) {
	h.logger.Println("Handling tools/list request")
	
	tools := h.server.registry.ListTools()
	result := mcp.ListToolsResult{
		Tools: tools,
	}
	
	return mcp.CreateResponse(msg.ID, result), nil
}

// handleToolsCall handles the tools/call request
func (h *Handler) handleToolsCall(msg *mcp.JSONRPCMessage) (*mcp.JSONRPCMessage, error) {
	h.logger.Println("Handling tools/call request")
	
	// Parse the call request
	var callReq mcp.CallToolRequest
	paramsData, err := json.Marshal(msg.Params)
	if err != nil {
		return mcp.CreateErrorResponse(msg.ID, mcp.NewInvalidParamsError("invalid params format")), nil
	}
	
	if err := json.Unmarshal(paramsData, &callReq); err != nil {
		return mcp.CreateErrorResponse(msg.ID, mcp.NewInvalidParamsError(err.Error())), nil
	}
	
	// Execute the tool
	result, err := h.server.registry.ExecuteTool(callReq.Name, callReq.Arguments)
	if err != nil {
		h.logger.Printf("Tool execution error: %v", err)
		// Return error as tool result
		errorResult := mcp.CallToolResult{
			Content: []mcp.Content{
				{
					Type: "text",
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
			IsError: true,
		}
		return mcp.CreateResponse(msg.ID, errorResult), nil
	}
	
	return mcp.CreateResponse(msg.ID, result), nil
}

// handlePing handles the ping request
func (h *Handler) handlePing(msg *mcp.JSONRPCMessage) (*mcp.JSONRPCMessage, error) {
	return mcp.CreateResponse(msg.ID, map[string]string{"status": "ok"}), nil
}

