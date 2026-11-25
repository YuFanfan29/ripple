// Command media processes uploaded images.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YuFanfan29/ripple/internal/env"
	"github.com/YuFanfan29/ripple/internal/logging"
	"github.com/YuFanfan29/ripple/internal/server"
)

func main() {
	logger := logging.New("media", os.Stdout)
	if err := run(logger); err != nil {
		logger.Error("media stopped with an error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(server.Config{
		Addr: env.String("MEDIA_HTTP_ADDR", ":8082"),
		// Draining and shutdown together stay under Docker's default
		// 10-second stop timeout, so the process exits before SIGKILL.
		DrainDelay:      2 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}, logger)
	return srv.Run(ctx)
}
