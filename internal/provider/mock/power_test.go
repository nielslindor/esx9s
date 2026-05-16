package mock

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/nielslindor/esx9s/internal/audit"
	"github.com/nielslindor/esx9s/internal/confirm"
	"github.com/nielslindor/esx9s/internal/domain"
)

func TestPowerServicePlansTypedConfirmationForDisruptiveAction(t *testing.T) {
	service := NewPowerService(samplePowerInventory(), nil)

	plan, err := service.Plan(context.Background(), "vm-on", PowerActionOff)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	if plan.Confirm.Level != confirm.LevelTyped {
		t.Fatalf("confirmation level = %v, want typed", plan.Confirm.Level)
	}
	if got, want := plan.Confirm.Token(), "POWER_OFF router"; got != want {
		t.Fatalf("confirmation token = %q, want %q", got, want)
	}
	if plan.DesiredState != domain.PowerStatePoweredOff {
		t.Fatalf("desired state = %q, want powered off", plan.DesiredState)
	}
}

func TestPowerServiceApplyMutatesMockInventoryAndWritesAudit(t *testing.T) {
	var buffer bytes.Buffer
	logger := audit.NewLogger(&buffer)
	service := NewPowerService(samplePowerInventory(), logger)
	service.now = func() time.Time {
		return time.Date(2026, 5, 16, 10, 30, 0, 0, time.UTC)
	}

	plan, err := service.Plan(context.Background(), "vm-off", PowerActionOn)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	task, err := service.Apply(context.Background(), plan, "yes", "local-user")
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if task.Status != domain.TaskStatusSuccess {
		t.Fatalf("task status = %q, want success", task.Status)
	}

	inventory, err := service.Inventory(context.Background())
	if err != nil {
		t.Fatalf("Inventory() error = %v", err)
	}
	vm, ok := findVM(inventory.VMs, "vm-off")
	if !ok {
		t.Fatal("vm-off not found after apply")
	}
	if vm.PowerState != domain.PowerStatePoweredOn {
		t.Fatalf("power state = %q, want powered_on", vm.PowerState)
	}
	if got, want := len(inventory.Tasks), 1; got != want {
		t.Fatalf("task count = %d, want %d", got, want)
	}

	var event audit.Event
	if err := json.Unmarshal([]byte(strings.TrimSpace(buffer.String())), &event); err != nil {
		t.Fatalf("unmarshal audit event: %v", err)
	}
	if event.Action != string(PowerActionOn) || event.Result != "success" || event.PlanID != plan.ID {
		t.Fatalf("audit event = %+v, want action/result/plan", event)
	}
}

func TestPowerServiceApplyRejectsMissingConfirmation(t *testing.T) {
	service := NewPowerService(samplePowerInventory(), nil)

	plan, err := service.Plan(context.Background(), "vm-on", PowerActionOff)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	_, err = service.Apply(context.Background(), plan, "yes", "local-user")
	if err == nil {
		t.Fatal("Apply() error = nil, want confirmation error")
	}
	if !strings.Contains(err.Error(), plan.Confirm.Token()) {
		t.Fatalf("Apply() error = %q, want confirmation token", err)
	}
}

func TestPowerServiceRejectsUnsupportedStateTransitions(t *testing.T) {
	service := NewPowerService(samplePowerInventory(), nil)

	_, err := service.Plan(context.Background(), "vm-off", PowerActionShutdownGuest)
	if err == nil {
		t.Fatal("Plan() error = nil, want state transition error")
	}
	if !strings.Contains(err.Error(), "powered-on") {
		t.Fatalf("Plan() error = %q, want powered-on", err)
	}
}

func samplePowerInventory() domain.Inventory {
	return domain.Inventory{
		VMs: []domain.VM{
			{
				ID:         "vm-on",
				Name:       "router",
				HostID:     "host-1",
				HostName:   "esxi01",
				PowerState: domain.PowerStatePoweredOn,
			},
			{
				ID:         "vm-off",
				Name:       "template",
				HostID:     "host-1",
				HostName:   "esxi01",
				PowerState: domain.PowerStatePoweredOff,
			},
		},
	}
}
