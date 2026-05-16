package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelRendersDefaultHostsView(t *testing.T) {
	model := NewModel()

	output := model.View()

	for _, want := range []string{"esx9s", "mode:mock", "Hosts", "esx-lab-01", "1/h hosts"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected default view to contain %q, got:\n%s", want, output)
		}
	}
}

func TestModelNavigatesBetweenViews(t *testing.T) {
	model := NewModel()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	model = updated.(Model)
	if got := model.ActiveView(); got != "VMs" {
		t.Fatalf("expected right navigation to select VMs, got %q", got)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("4")})
	model = updated.(Model)
	if got := model.ActiveView(); got != "Events/Audit" {
		t.Fatalf("expected number navigation to select Events/Audit, got %q", got)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")})
	model = updated.(Model)
	if got := model.ActiveView(); got != "VMs" {
		t.Fatalf("expected v navigation to select VMs, got %q", got)
	}
}

func TestModelMovesCursorWithinActiveView(t *testing.T) {
	model := NewModel()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	model = updated.(Model)

	if got := model.cursor[hostsView]; got != 1 {
		t.Fatalf("expected cursor to move down once, got %d", got)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	model = updated.(Model)

	if got := model.cursor[hostsView]; got != 0 {
		t.Fatalf("expected cursor to move back to top, got %d", got)
	}
}

func TestModelHelpAndPlaceholders(t *testing.T) {
	model := NewModel()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	model = updated.(Model)
	if output := model.View(); !strings.Contains(output, "Keyboard") {
		t.Fatalf("expected help to render, got:\n%s", output)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")})
	model = updated.(Model)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	model = updated.(Model)
	if output := model.View(); !strings.Contains(output, "typed confirmation and audit") {
		t.Fatalf("expected power placeholder status, got:\n%s", output)
	}
}
