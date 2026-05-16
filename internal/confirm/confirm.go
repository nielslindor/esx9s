package confirm

import (
	"errors"
	"fmt"
	"strings"
)

// Level describes how deliberate an operator confirmation must be.
type Level int

const (
	LevelNone Level = iota
	LevelVisible
	LevelTyped
)

// Request is the typed confirmation contract attached to an action plan.
type Request struct {
	Level      Level
	Action     string
	TargetType string
	TargetName string
}

// Token returns the exact phrase an operator must type for LevelTyped actions.
func (r Request) Token() string {
	action := strings.TrimSpace(r.Action)
	target := strings.TrimSpace(r.TargetName)
	if action == "" || target == "" {
		return ""
	}

	return fmt.Sprintf("%s %s", strings.ToUpper(action), target)
}

// Validate checks that the confirmation request can safely be shown to an operator.
func (r Request) Validate() error {
	var missing []string
	if strings.TrimSpace(r.Action) == "" {
		missing = append(missing, "action")
	}
	if strings.TrimSpace(r.TargetType) == "" {
		missing = append(missing, "target_type")
	}
	if strings.TrimSpace(r.TargetName) == "" {
		missing = append(missing, "target_name")
	}
	if r.Level < LevelNone || r.Level > LevelTyped {
		missing = append(missing, "level")
	}
	if len(missing) > 0 {
		return fmt.Errorf("confirmation request missing required field(s): %s", strings.Join(missing, ", "))
	}

	return nil
}

// Check returns nil when input satisfies the request's confirmation level.
func Check(request Request, input string) error {
	if err := request.Validate(); err != nil {
		return err
	}

	switch request.Level {
	case LevelNone:
		return nil
	case LevelVisible:
		if strings.EqualFold(strings.TrimSpace(input), "yes") {
			return nil
		}
		return errors.New(`confirmation requires "yes"`)
	case LevelTyped:
		if strings.TrimSpace(input) == request.Token() {
			return nil
		}
		return fmt.Errorf("confirmation requires exact phrase %q", request.Token())
	default:
		return errors.New("unsupported confirmation level")
	}
}
