package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
)

// Transport defines the interface for message transport
type Transport interface {
	// Start starts the transport
	Start(ctx context.Context) error
	
	// Send sends a message
	Send(msg *mcp.JSONRPCMessage) error
	
	// Receive receives messages and passes them to the handler
	Receive(handler MessageHandler) error
	
	// Close closes the transport
	Close() error
}

// MessageHandler is a function that handles incoming messages
type MessageHandler func(msg *mcp.JSONRPCMessage) (*mcp.JSONRPCMessage, error)

// StdioTransport implements stdio transport
type StdioTransport struct {
	reader *bufio.Reader
	writer io.Writer
	logger *log.Logger
}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport() *StdioTransport {
	return &StdioTransport{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
		logger: log.New(os.Stderr, "[Transport] ", log.LstdFlags),
	}
}

// Start starts the stdio transport
func (t *StdioTransport) Start(ctx context.Context) error {
	t.logger.Println("Stdio transport started")
	return nil
}

// Send sends a message through stdio
func (t *StdioTransport) Send(msg *mcp.JSONRPCMessage) error {
	data, err := mcp.EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}
	
	// Write the message followed by newline
	if _, err := t.writer.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	
	t.logger.Printf("Sent message: method=%s, id=%v", msg.Method, msg.ID)
	return nil
}

// Receive receives and processes messages from stdio
func (t *StdioTransport) Receive(handler MessageHandler) error {
	t.logger.Println("Starting to receive messages...")
	
	for {
		// Read a line
		line, err := t.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				t.logger.Println("EOF reached, closing transport")
				return nil
			}
			return fmt.Errorf("failed to read line: %w", err)
		}
		
		// Skip empty lines
		if len(line) <= 1 {
			continue
		}
		
		// Parse the message
		msg, err := mcp.ParseMessage(line)
		if err != nil {
			t.logger.Printf("Failed to parse message: %v", err)
			// Send error response
			errMsg := mcp.CreateErrorResponse(nil, mcp.NewParseError(err.Error()))
			if err := t.Send(errMsg); err != nil {
				t.logger.Printf("Failed to send error response: %v", err)
			}
			continue
		}
		
		t.logger.Printf("Received message: method=%s, id=%v", msg.Method, msg.ID)
		
		// Handle the message
		response, err := handler(msg)
		if err != nil {
			t.logger.Printf("Handler error: %v", err)
			// Send error response
			errMsg := mcp.CreateErrorResponse(msg.ID, mcp.NewInternalError(err.Error()))
			if err := t.Send(errMsg); err != nil {
				t.logger.Printf("Failed to send error response: %v", err)
			}
			continue
		}
		
		// Send response (if not a notification)
		if response != nil && msg.ID != nil {
			if err := t.Send(response); err != nil {
				t.logger.Printf("Failed to send response: %v", err)
			}
		}
	}
}

// Close closes the stdio transport
func (t *StdioTransport) Close() error {
	t.logger.Println("Stdio transport closed")
	return nil
}

