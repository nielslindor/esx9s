package audit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const redacted = "[REDACTED]"

// Event is the append-only audit record written as one JSON object per line.
type Event struct {
	Timestamp  time.Time `json:"timestamp"`
	Operator   string    `json:"operator"`
	TargetHost string    `json:"target_host"`
	TargetType string    `json:"target_type"`
	TargetID   string    `json:"target_id"`
	TargetName string    `json:"target_name"`
	Action     string    `json:"action"`
	PlanID     string    `json:"plan_id"`
	Result     string    `json:"result"`
	Error      *string   `json:"error"`
}

// Validate rejects incomplete or unsafe audit events before they reach disk.
func (e Event) Validate() error {
	var missing []string

	if e.Timestamp.IsZero() {
		missing = append(missing, "timestamp")
	}
	if strings.TrimSpace(e.Operator) == "" {
		missing = append(missing, "operator")
	}
	if strings.TrimSpace(e.TargetHost) == "" {
		missing = append(missing, "target_host")
	}
	if strings.TrimSpace(e.TargetType) == "" {
		missing = append(missing, "target_type")
	}
	if strings.TrimSpace(e.TargetID) == "" {
		missing = append(missing, "target_id")
	}
	if strings.TrimSpace(e.Action) == "" {
		missing = append(missing, "action")
	}
	if strings.TrimSpace(e.PlanID) == "" {
		missing = append(missing, "plan_id")
	}
	if strings.TrimSpace(e.Result) == "" {
		missing = append(missing, "result")
	}

	if len(missing) > 0 {
		return fmt.Errorf("audit event missing required field(s): %s", strings.Join(missing, ", "))
	}

	if containsInlineSecret(e.Operator, e.TargetHost, e.TargetType, e.TargetID, e.TargetName, e.Action, e.PlanID, e.Result) {
		return errors.New("audit event contains secret material in a structured field")
	}

	return nil
}

// Redacted returns a copy safe to persist in an audit log.
func (e Event) Redacted() Event {
	if e.Error != nil {
		errText := RedactString(*e.Error)
		e.Error = &errText
	}

	return e
}

// Logger appends audit events to a JSONL stream.
type Logger struct {
	mu     sync.Mutex
	writer io.Writer
	closer io.Closer
}

// NewLogger writes JSONL audit events to writer. It is useful for tests and
// callers that manage their own file lifecycle.
func NewLogger(writer io.Writer) *Logger {
	return &Logger{writer: writer}
}

// Open opens path for append-only JSONL audit logging, creating it when needed.
func Open(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}

	return &Logger{writer: file, closer: file}, nil
}

// Append validates, redacts, and writes event as a single JSONL record.
func (l *Logger) Append(event Event) error {
	if l == nil || l.writer == nil {
		return errors.New("audit logger has no writer")
	}

	event = event.Redacted()
	if err := event.Validate(); err != nil {
		return err
	}

	line, err := json.Marshal(event)
	if err != nil {
		return err
	}
	line = append(line, '\n')

	l.mu.Lock()
	defer l.mu.Unlock()

	n, err := l.writer.Write(line)
	if err == nil && n != len(line) {
		return io.ErrShortWrite
	}
	return err
}

// Close releases the underlying file when the logger owns one.
func (l *Logger) Close() error {
	if l == nil || l.closer == nil {
		return nil
	}

	return l.closer.Close()
}
