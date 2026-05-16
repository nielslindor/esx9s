package confirm

import (
	"strings"
	"testing"
)

func TestTypedTokenIncludesUppercaseActionAndTarget(t *testing.T) {
	request := Request{
		Level:      LevelTyped,
		Action:     "power_off",
		TargetType: "vm",
		TargetName: "home-assistant",
	}

	if got, want := request.Token(), "POWER_OFF home-assistant"; got != want {
		t.Fatalf("Token() = %q, want %q", got, want)
	}
}

func TestCheckAcceptsVisibleConfirmation(t *testing.T) {
	request := Request{
		Level:      LevelVisible,
		Action:     "power_on",
		TargetType: "vm",
		TargetName: "template-debian12",
	}

	if err := Check(request, " yes "); err != nil {
		t.Fatalf("Check() error = %v", err)
	}
}

func TestCheckRejectsTypedConfirmationMismatch(t *testing.T) {
	request := Request{
		Level:      LevelTyped,
		Action:     "power_off",
		TargetType: "vm",
		TargetName: "home-assistant",
	}

	err := Check(request, "yes")
	if err == nil {
		t.Fatal("Check() error = nil, want mismatch")
	}
	if !strings.Contains(err.Error(), request.Token()) {
		t.Fatalf("Check() error = %q, want token %q", err, request.Token())
	}
}

func TestCheckRejectsIncompleteRequest(t *testing.T) {
	err := Check(Request{Level: LevelVisible, Action: "power_on"}, "yes")
	if err == nil {
		t.Fatal("Check() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "target_type") {
		t.Fatalf("Check() error = %q, want target_type", err)
	}
}
