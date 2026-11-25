package server_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/YuFanfan29/ripple/internal/server"
)

func newServer(cfg server.Config) *server.Server {
	return server.New(cfg, slog.New(slog.DiscardHandler))
}

func TestLivezReportsOK(t *testing.T) {
	srv := newServer(server.Config{})
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/livez", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /livez = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyzIsUnavailableBeforeServing(t *testing.T) {
	srv := newServer(server.Config{})
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("GET /readyz = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestServeIsReadyUntilContextIsCanceled(t *testing.T) {
	srv := newServer(server.Config{ShutdownTimeout: time.Second})
	base, stop, done := start(t, srv)

	waitForStatus(t, base+"/readyz", http.StatusOK)
	stop()

	waitForResult(t, done)
}

func TestServeFinishesInFlightRequestsBeforeReturning(t *testing.T) {
	srv := newServer(server.Config{ShutdownTimeout: 5 * time.Second})
	started := make(chan struct{})
	release := make(chan struct{})
	srv.Handle("GET /slow", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	base, stop, done := start(t, srv)
	waitForStatus(t, base+"/readyz", http.StatusOK)

	result := make(chan error, 1)
	go func() {
		code, err := get(t.Context(), base+"/slow")
		if err == nil && code != http.StatusOK {
			err = fmt.Errorf("status %d", code)
		}
		result <- err
	}()
	<-started
	stop()
	waitForRefused(t, base+"/livez")
	close(release)

	if err := <-result; err != nil {
		t.Fatalf("in-flight request failed: %v", err)
	}
	waitForResult(t, done)
}

func TestReadyzIsUnavailableWhileDraining(t *testing.T) {
	srv := newServer(server.Config{DrainDelay: time.Second, ShutdownTimeout: time.Second})
	base, stop, done := start(t, srv)
	waitForStatus(t, base+"/readyz", http.StatusOK)

	stop()

	waitForStatus(t, base+"/readyz", http.StatusServiceUnavailable)
	if code, err := get(t.Context(), base+"/livez"); err != nil || code != http.StatusOK {
		t.Fatalf("GET /livez while draining = %d (error: %v), want %d", code, err, http.StatusOK)
	}
	waitForResult(t, done)
}

func TestServeLogsLifecycleEvents(t *testing.T) {
	var logs bytes.Buffer
	srv := server.New(server.Config{}, slog.New(slog.NewJSONHandler(&logs, nil)))
	base, stop, done := start(t, srv)
	waitForStatus(t, base+"/readyz", http.StatusOK)

	stop()
	waitForResult(t, done)

	for _, msg := range []string{"server started", "draining", "server stopped"} {
		if !strings.Contains(logs.String(), `"msg":"`+msg+`"`) {
			t.Errorf("logs do not contain %q:\n%s", msg, logs.String())
		}
	}
}

func TestRunFailsWhenAddressIsInUse(t *testing.T) {
	var lc net.ListenConfig
	taken, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = taken.Close() })
	srv := newServer(server.Config{Addr: taken.Addr().String()})

	if err := srv.Run(t.Context()); err == nil {
		t.Fatal("Run returned nil, want an error for an address that is in use")
	}
}

func TestServeGivesUpAfterShutdownTimeout(t *testing.T) {
	srv := newServer(server.Config{ShutdownTimeout: 100 * time.Millisecond})
	started := make(chan struct{})
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	srv.Handle("GET /stuck", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
	}))
	base, stop, done := start(t, srv)
	waitForStatus(t, base+"/readyz", http.StatusOK)
	go func() { _, _ = get(t.Context(), base+"/stuck") }()
	<-started

	stop()

	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Serve returned %v, want context.DeadlineExceeded", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not give up after ShutdownTimeout")
	}
}

// start runs srv.Serve on a free local port. It returns the base URL, a
// function that cancels the server's context, and a channel that receives
// the value Serve returns.
func start(t *testing.T, srv *server.Server) (string, context.CancelFunc, <-chan error) {
	t.Helper()
	var lc net.ListenConfig
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	done := make(chan error, 1)
	go func() { done <- srv.Serve(ctx, ln) }()
	return "http://" + ln.Addr().String(), cancel, done
}

var client = &http.Client{Timeout: time.Second}

// get sends a GET request to url and returns the response status code.
func get(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

// waitForStatus polls url until it returns want, and fails the test after
// two seconds.
func waitForStatus(t *testing.T, url string, want int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		got, err := get(t.Context(), url)
		if err == nil && got == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("GET %s = %d (error: %v), want %d", url, got, err, want)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// waitForRefused polls url until the server stops accepting connections,
// and fails the test after two seconds.
func waitForRefused(t *testing.T, url string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err := get(t.Context(), url); err != nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("GET %s still succeeds, want the connection to be refused", url)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// waitForResult waits for Serve to return, and fails the test if Serve
// returns an error or takes longer than five seconds.
func waitForResult(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Serve returned %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not return after its context was canceled")
	}
}
