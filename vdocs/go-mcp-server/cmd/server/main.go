package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/server"
)

var (
	transport             string
	host                  string
	port                  int
	services              string
	headersCredentialOnly bool
	env                   string
)

func init() {
	flag.StringVar(&transport, "transport", "stdio", "Transport type: stdio, http, sse")
	flag.StringVar(&host, "host", "127.0.0.1", "Host address for HTTP/SSE transport")
	flag.IntVar(&port, "port", 8000, "Port number for HTTP/SSE transport")
	flag.StringVar(&services, "services", "", "Comma-separated list of services to enable (e.g., ecs,vpc,rds)")
	flag.BoolVar(&headersCredentialOnly, "headers-credential-only", false, "Use credentials only from HTTP headers")
	flag.StringVar(&env, "env", "domestic", "Environment: domestic or international")
}

func main() {
	flag.Parse()
	
	// Create configuration
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Override with command line flags
	if transport != "" {
		cfg.Transport = transport
	}
	if host != "" {
		cfg.Host = host
	}
	if port != 0 {
		cfg.Port = port
	}
	if services != "" {
		cfg.Services = parseServices(services)
	}
	if headersCredentialOnly {
		cfg.HeadersCredentialOnly = headersCredentialOnly
	}
	if env != "" {
		cfg.Environment = env
	}
	
	// Apply settings
	settings := config.GetSettings()
	settings.SetHeadersCredentialOnly(cfg.HeadersCredentialOnly)
	if err := settings.SetEnvironment(cfg.Environment); err != nil {
		log.Fatalf("Invalid environment: %v", err)
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	
	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(ctx); err != nil {
			errChan <- err
		}
	}()
	
	// Wait for signal or error
	select {
	case <-sigChan:
		log.Println("Received interrupt signal, shutting down...")
		cancel()
		if err := srv.Stop(); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	}
}

func parseServices(s string) []string {
	parts := strings.Split(s, ",")
	services := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			services = append(services, strings.ToLower(trimmed))
		}
	}
	return services
}

