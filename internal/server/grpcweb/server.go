package grpcweb

import (
	"context"
	"net/http"
	"time"

	"vorpalstacks/internal/core/logs"
)

// Server is a gRPC-Web server that can handle both HTTP and WebSocket connections.
type Server struct {
	config     *Config
	httpServer *http.Server
	mux        *http.ServeMux
}

// NewServer creates a new gRPC-Web server with the given configuration.
func NewServer(cfg *Config) *Server {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Server{
		config: cfg,
		httpServer: &http.Server{
			Addr:              cfg.getBindAddr() + ":" + cfg.DefaultPort(),
			ReadHeaderTimeout: 10 * time.Second,
		},
		mux: http.NewServeMux(),
	}
}

// Handle registers an HTTP handler for the given pattern.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// HandleFunc registers an HTTP handler function for the given pattern.
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

// Start starts the gRPC-Web server and blocks until it stops.
func (s *Server) Start() error {
	s.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logs.Error("PANIC in gRPC-Web handler", logs.Any("panic", rec))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		s.mux.ServeHTTP(w, r)
	})
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Addr returns the server's listening address.
func (s *Server) Addr() string {
	return s.httpServer.Addr
}

// Port returns the server's port number.
func (s *Server) Port() string {
	return s.config.DefaultPort()
}

// URL returns the server's base URL.
func (s *Server) URL() string {
	return "http://" + s.httpServer.Addr
}
