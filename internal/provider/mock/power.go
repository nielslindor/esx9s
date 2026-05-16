package mock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nielslindor/esx9s/internal/audit"
	"github.com/nielslindor/esx9s/internal/confirm"
	"github.com/nielslindor/esx9s/internal/domain"
)

// PowerAction is a mock-only VM power action. It does not call ESXi.
type PowerAction string

const (
	PowerActionOn            PowerAction = "power_on"
	PowerActionShutdownGuest PowerAction = "shutdown_guest"
	PowerActionOff           PowerAction = "power_off"
	PowerActionReset         PowerAction = "reset"
)

// AuditAppender is the audit package integration point used by the mock action service.
type AuditAppender interface {
	Append(audit.Event) error
}

// PowerService simulates VM power actions against in-memory mock inventory.
type PowerService struct {
	mu        sync.Mutex
	inventory domain.Inventory
	audit     AuditAppender
	now       func() time.Time
	nextTask  int
}

// PowerPlan describes the visible plan an operator must confirm before apply.
type PowerPlan struct {
	ID           string
	Action       PowerAction
	VMID         string
	VMName       string
	HostID       string
	HostName     string
	CurrentState domain.PowerState
	DesiredState domain.PowerState
	Confirm      confirm.Request
}

// NewPowerService returns a mock power action service seeded with inventory.
func NewPowerService(inventory domain.Inventory, auditLog AuditAppender) *PowerService {
	return &PowerService{
		inventory: cloneInventory(inventory),
		audit:     auditLog,
		now:       time.Now,
		nextTask:  1,
	}
}

// Plan builds a mock action plan without changing inventory.
func (s *PowerService) Plan(ctx context.Context, vmID string, action PowerAction) (PowerPlan, error) {
	if err := ctx.Err(); err != nil {
		return PowerPlan{}, err
	}
	if s == nil {
		return PowerPlan{}, errors.New("mock power service is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	vm, ok := findVM(s.inventory.VMs, vmID)
	if !ok {
		return PowerPlan{}, fmt.Errorf("vm %q not found", vmID)
	}

	desired, level, err := desiredPowerState(vm.PowerState, action)
	if err != nil {
		return PowerPlan{}, err
	}

	return PowerPlan{
		ID:           fmt.Sprintf("mock-%s-%s", action, vm.ID),
		Action:       action,
		VMID:         vm.ID,
		VMName:       vm.Name,
		HostID:       vm.HostID,
		HostName:     vm.HostName,
		CurrentState: vm.PowerState,
		DesiredState: desired,
		Confirm: confirm.Request{
			Level:      level,
			Action:     string(action),
			TargetType: "vm",
			TargetName: vm.Name,
		},
	}, nil
}

// Apply validates confirmation, mutates mock inventory, records a mock task, and writes audit when configured.
func (s *PowerService) Apply(ctx context.Context, plan PowerPlan, confirmationInput string, operator string) (domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return domain.Task{}, err
	}
	if s == nil {
		return domain.Task{}, errors.New("mock power service is nil")
	}
	if err := confirm.Check(plan.Confirm, confirmationInput); err != nil {
		return domain.Task{}, err
	}

	now := s.now().UTC()
	auditLog := s.audit
	auditEvent := audit.Event{}
	if auditLog != nil {
		auditEvent = audit.Event{
			Timestamp:  now,
			Operator:   operator,
			TargetHost: plan.HostName,
			TargetType: "vm",
			TargetID:   plan.VMID,
			TargetName: plan.VMName,
			Action:     string(plan.Action),
			PlanID:     plan.ID,
			Result:     "success",
		}
		if err := auditEvent.Validate(); err != nil {
			return domain.Task{}, err
		}
	}

	s.mu.Lock()
	vmIndex, ok := findVMIndex(s.inventory.VMs, plan.VMID)
	if !ok {
		s.mu.Unlock()
		return domain.Task{}, fmt.Errorf("vm %q not found", plan.VMID)
	}
	if s.inventory.VMs[vmIndex].PowerState != plan.CurrentState {
		current := s.inventory.VMs[vmIndex].PowerState
		s.mu.Unlock()
		return domain.Task{}, fmt.Errorf("vm %q state changed from %s to %s", plan.VMName, plan.CurrentState, current)
	}

	finished := now
	task := domain.Task{
		ID:         fmt.Sprintf("mock-task-%04d", s.nextTask),
		HostID:     plan.HostID,
		HostName:   plan.HostName,
		TargetID:   plan.VMID,
		TargetName: plan.VMName,
		Operation:  string(plan.Action),
		Status:     domain.TaskStatusSuccess,
		StartedAt:  now,
		FinishedAt: &finished,
	}
	s.nextTask++
	s.inventory.VMs[vmIndex].PowerState = plan.DesiredState
	s.inventory.Tasks = append(s.inventory.Tasks, task)
	s.mu.Unlock()

	if auditLog != nil {
		if err := auditLog.Append(auditEvent); err != nil {
			return domain.Task{}, err
		}
	}

	return task, nil
}

// Inventory returns a defensive copy of the mock action service inventory.
func (s *PowerService) Inventory(ctx context.Context) (domain.Inventory, error) {
	if err := ctx.Err(); err != nil {
		return domain.Inventory{}, err
	}
	if s == nil {
		return domain.Inventory{}, errors.New("mock power service is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return cloneInventory(s.inventory), nil
}

func desiredPowerState(current domain.PowerState, action PowerAction) (domain.PowerState, confirm.Level, error) {
	switch action {
	case PowerActionOn:
		if current == domain.PowerStatePoweredOn {
			return "", 0, errors.New("vm is already powered on")
		}
		return domain.PowerStatePoweredOn, confirm.LevelVisible, nil
	case PowerActionShutdownGuest:
		if current != domain.PowerStatePoweredOn {
			return "", 0, errors.New("guest shutdown requires a powered-on vm")
		}
		return domain.PowerStatePoweredOff, confirm.LevelVisible, nil
	case PowerActionOff:
		if current == domain.PowerStatePoweredOff {
			return "", 0, errors.New("vm is already powered off")
		}
		return domain.PowerStatePoweredOff, confirm.LevelTyped, nil
	case PowerActionReset:
		if current != domain.PowerStatePoweredOn {
			return "", 0, errors.New("reset requires a powered-on vm")
		}
		return domain.PowerStatePoweredOn, confirm.LevelTyped, nil
	default:
		return "", 0, fmt.Errorf("unsupported power action %q", action)
	}
}

func findVM(vms []domain.VM, id string) (domain.VM, bool) {
	for _, vm := range vms {
		if vm.ID == id {
			return vm, true
		}
	}
	return domain.VM{}, false
}

func findVMIndex(vms []domain.VM, id string) (int, bool) {
	for i, vm := range vms {
		if vm.ID == id {
			return i, true
		}
	}
	return -1, false
}
