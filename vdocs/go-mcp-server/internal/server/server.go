package server

import (
	"context"
	"fmt"
	"log"
	"os"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/tools"
)

// Server represents the MCP server
type Server struct {
	config    *config.Config
	transport Transport
	registry  *tools.Registry
	handler   *Handler
	logger    *log.Logger
}

// NewServer creates a new MCP server
func NewServer(cfg *config.Config) (*Server, error) {
	logger := log.New(os.Stderr, "[Server] ", log.LstdFlags)
	
	// Create tool registry
	registry := tools.NewRegistry(cfg, logger)
	
	server := &Server{
		config:   cfg,
		registry: registry,
		logger:   logger,
	}
	
	// Create handler
	server.handler = NewHandler(server)
	
	// Create transport based on config
	switch cfg.Transport {
	case "stdio":
		server.transport = NewStdioTransport()
	case "http":
		return nil, fmt.Errorf("HTTP transport not yet implemented")
	case "sse":
		return nil, fmt.Errorf("SSE transport not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported transport: %s", cfg.Transport)
	}
	
	return server, nil
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	s.logger.Printf("Starting MCP server (transport: %s)", s.config.Transport)
	
	// Start transport
	if err := s.transport.Start(ctx); err != nil {
		return fmt.Errorf("failed to start transport: %w", err)
	}
	
	// Register tools
	if err := s.registry.RegisterAllTools(); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}
	
	s.logger.Printf("Server started with %d tools registered", len(s.registry.ListTools()))
	
	// Start receiving messages
	return s.transport.Receive(s.handler.Handle)
}

// Stop stops the server
func (s *Server) Stop() error {
	s.logger.Println("Stopping MCP server")
	return s.transport.Close()
}

// GetRegistry returns the tool registry
func (s *Server) GetRegistry() *tools.Registry {
	return s.registry
}

