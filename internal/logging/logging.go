// Package logging builds the structured JSON logger that every ripple Go
// service uses.
package logging

import (
	"io"
	"log/slog"
)

// New returns a logger that writes one JSON object per line to w. Every line
// carries the service name, and the message key is "message", matching the
// log fields that all ripple services share.
func New(service string, w io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	})
	return slog.New(handler).With("service", service)
}
