package logging_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/YuFanfan29/ripple/internal/logging"
)

func TestNewWritesConventionFields(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New("feed", &buf)

	logger.Info("server started", "addr", ":8081")

	var line map[string]any
	if err := json.Unmarshal(buf.Bytes(), &line); err != nil {
		t.Fatalf("log line is not JSON: %v\n%s", err, buf.String())
	}
	for key, want := range map[string]any{
		"level":   "INFO",
		"service": "feed",
		"message": "server started",
		"addr":    ":8081",
	} {
		if got := line[key]; got != want {
			t.Errorf("%s = %v, want %v", key, got, want)
		}
	}
	if _, ok := line["time"]; !ok {
		t.Error("time field is missing")
	}
	if _, ok := line["msg"]; ok {
		t.Error(`"msg" should be renamed to "message"`)
	}
}
