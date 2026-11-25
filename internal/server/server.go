// Package server runs the HTTP server shared by ripple's Go services. It
// serves liveness and readiness checks and shuts down gracefully, so a deploy
// can drain traffic before a process exits.
package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

// Config controls how a Server listens and shuts down.
type Config struct {
	// Addr is the TCP address to listen on, such as ":8081".
	Addr string
	// DrainDelay is how long the server keeps serving after readiness turns
	// false, so load balancers stop sending new requests before it closes.
	DrainDelay time.Duration
	// ShutdownTimeout is how long the server waits for in-flight requests
	// after draining. Serve returns an error if they take longer.
	ShutdownTimeout time.Duration
}

// Server is an HTTP server with liveness and readiness checks and graceful
// shutdown.
type Server struct {
	cfg    Config
	logger *slog.Logger
	mux    *http.ServeMux
	ready  atomic.Bool
}

// New returns a Server that serves GET /livez and GET /readyz. Register
// application routes with Handle before calling Run or Serve.
func New(cfg Config, logger *slog.Logger) *Server {
	s := &Server{cfg: cfg, logger: logger, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /livez", livez)
	s.mux.HandleFunc("GET /readyz", s.readyz)
	return s
}

// Handle registers an application handler for pattern.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// Run listens on the configured address and serves until ctx is canceled.
// See Serve for how it shuts down.
func (s *Server) Run(ctx context.Context) error {
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.cfg.Addr, err)
	}
	return s.Serve(ctx, ln)
}

// Serve accepts connections on ln until ctx is canceled. Serve takes
// ownership of ln.
//
// While serving, /readyz reports ready. When ctx is canceled, the server:
//
//  1. reports not ready, and keeps serving for Config.DrainDelay;
//  2. stops accepting connections and waits up to Config.ShutdownTimeout
//     for in-flight requests to finish.
//
// Serve returns nil after a clean shutdown.
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	srv := &http.Server{Handler: s.mux, ReadHeaderTimeout: 5 * time.Second}
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()
	s.ready.Store(true)
	s.logger.Info("server started", "addr", ln.Addr().String())

	select {
	case err := <-serveErr:
		s.ready.Store(false)
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	s.ready.Store(false)
	s.logger.Info("draining", "delay", s.cfg.DrainDelay.String())
	drain := time.NewTimer(s.cfg.DrainDelay)
	defer drain.Stop()
	select {
	case <-drain.C:
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down: %w", err)
	}
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve: %w", err)
	}
	s.logger.Info("server stopped")
	return nil
}

// ServeHTTP lets the Server be used directly as an http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func livez(w http.ResponseWriter, _ *http.Request) {
	writeStatus(w, http.StatusOK)
}

func (s *Server) readyz(w http.ResponseWriter, _ *http.Request) {
	if !s.ready.Load() {
		writeStatus(w, http.StatusServiceUnavailable)
		return
	}
	writeStatus(w, http.StatusOK)
}

func writeStatus(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	_, _ = io.WriteString(w, http.StatusText(code)+"\n")
}
