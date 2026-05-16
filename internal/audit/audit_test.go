package audit

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoggerAppendWritesJSONLine(t *testing.T) {
	var buffer bytes.Buffer
	logger := NewLogger(&buffer)

	errText := "task failed: password=hunter2 token: abc123 Authorization: Bearer xyz"
	event := validEvent()
	event.Error = &errText

	if err := logger.Append(event); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(buffer.String(), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("got %d JSONL lines, want 1: %q", len(lines), buffer.String())
	}

	var got Event
	if err := json.Unmarshal([]byte(lines[0]), &got); err != nil {
		t.Fatalf("unmarshal audit event: %v", err)
	}

	if got.Action != event.Action {
		t.Fatalf("action = %q, want %q", got.Action, event.Action)
	}
	if got.Error == nil {
		t.Fatal("error was nil, want redacted error text")
	}
	if strings.Contains(*got.Error, "hunter2") || strings.Contains(*got.Error, "abc123") || strings.Contains(*got.Error, "xyz") {
		t.Fatalf("error was not redacted: %q", *got.Error)
	}
	if !strings.Contains(*got.Error, redacted) {
		t.Fatalf("error does not contain redaction marker: %q", *got.Error)
	}
}

func TestLoggerAppendRejectsMissingRequiredFields(t *testing.T) {
	var buffer bytes.Buffer
	logger := NewLogger(&buffer)

	event := validEvent()
	event.TargetID = ""

	err := logger.Append(event)
	if err == nil {
		t.Fatal("Append() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "target_id") {
		t.Fatalf("validation error = %q, want target_id", err)
	}
	if buffer.Len() != 0 {
		t.Fatalf("buffer length = %d, want 0", buffer.Len())
	}
}

func TestLoggerAppendRejectsStructuredSecretMaterial(t *testing.T) {
	var buffer bytes.Buffer
	logger := NewLogger(&buffer)

	event := validEvent()
	event.TargetHost = "https://operator:password@example.local"

	err := logger.Append(event)
	if err == nil {
		t.Fatal("Append() error = nil, want validation error")
	}
	if buffer.Len() != 0 {
		t.Fatalf("buffer length = %d, want 0", buffer.Len())
	}
}

func TestOpenAppendsToExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	if err := os.WriteFile(path, []byte("{\"existing\":true}\n"), 0o600); err != nil {
		t.Fatalf("seed audit log: %v", err)
	}

	logger, err := Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := logger.Append(validEvent()); err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	if err := logger.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read audit log: %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), string(content))
	}
	if lines[0] != "{\"existing\":true}" {
		t.Fatalf("first line = %q, want existing line preserved", lines[0])
	}
}

func TestOpenCreatesPrivateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")

	logger, err := Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := logger.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat audit log: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("audit log permissions = %v, want 0600", got)
	}
}

func TestOpenCreatesParentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "esx9s", "audit.jsonl")

	logger, err := Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := logger.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat audit log: %v", err)
	}
}

func TestRedactStringRedactsCommonSecretForms(t *testing.T) {
	input := "password=abc secret: def token='ghi' Bearer jkl https://user:pass@example.local/sdk"
	got := RedactString(input)

	for _, leaked := range []string{"abc", "def", "ghi", "jkl", "user:pass"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("RedactString(%q) leaked %q as %q", input, leaked, got)
		}
	}
}

func TestRedactFieldsRedactsSecretKeysAndValues(t *testing.T) {
	got := RedactFields(map[string]string{
		"action":   "power_off",
		"password": "hunter2",
		"message":  "failed with token=abc123",
	})

	if got["action"] != "power_off" {
		t.Fatalf("action = %q, want power_off", got["action"])
	}
	if got["password"] != redacted {
		t.Fatalf("password = %q, want redacted", got["password"])
	}
	if strings.Contains(got["message"], "abc123") {
		t.Fatalf("message leaked token: %q", got["message"])
	}
}

func validEvent() Event {
	return Event{
		Timestamp:  time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC),
		Operator:   "local-user",
		TargetHost: "esxi05",
		TargetType: "vm",
		TargetID:   "vm-123",
		TargetName: "opnsense",
		Action:     "power_off",
		PlanID:     "plan-123",
		Result:     "success",
	}
}
