package actions

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/nielslindor/esx9s/internal/domain"
	"github.com/nielslindor/esx9s/internal/provider"
)

const targetTypeVM = "vm"

// SnapshotCreateRequest contains the operator intent needed to plan a snapshot create.
type SnapshotCreateRequest struct {
	Operator      string
	VMID          string
	Name          string
	Description   string
	IncludeMemory bool
	Quiesce       bool
}

// SnapshotDeleteRequest contains the operator intent needed to plan a snapshot delete.
type SnapshotDeleteRequest struct {
	Operator       string
	VMID           string
	SnapshotID     string
	SnapshotName   string
	DeleteChildren bool
}

// Planner builds safe, auditable action plans from read-only provider inventory.
type Planner struct {
	Provider provider.Provider
}

func (p Planner) PlanSnapshotCreate(ctx context.Context, request SnapshotCreateRequest) (domain.ActionPlan, error) {
	inventory, err := p.inventory(ctx)
	if err != nil {
		return domain.ActionPlan{}, err
	}

	vm, ok := findVM(inventory.VMs, request.VMID)
	if !ok {
		return domain.ActionPlan{}, fmt.Errorf("VM %q not found", request.VMID)
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		return domain.ActionPlan{}, errors.New("snapshot name is required")
	}
	operator := strings.TrimSpace(request.Operator)
	if operator == "" {
		return domain.ActionPlan{}, errors.New("operator is required")
	}

	plan := domain.ActionPlan{
		Action:            domain.ActionCreateSnapshot,
		ConfirmationLevel: domain.ConfirmationVisibleTarget,
		ConfirmationToken: "yes",
		Operator:          operator,
		TargetHost:        vm.HostName,
		TargetType:        targetTypeVM,
		TargetID:          vm.ID,
		TargetName:        vm.Name,
		SnapshotName:      name,
		SnapshotDesc:      strings.TrimSpace(request.Description),
		IncludeMemory:     request.IncludeMemory,
		Quiesce:           request.Quiesce,
	}
	plan.ID = planID(plan, plan.SnapshotName, plan.SnapshotDesc)

	return plan, nil
}

func (p Planner) PlanSnapshotDelete(ctx context.Context, request SnapshotDeleteRequest) (domain.ActionPlan, error) {
	inventory, err := p.inventory(ctx)
	if err != nil {
		return domain.ActionPlan{}, err
	}

	vm, ok := findVM(inventory.VMs, request.VMID)
	if !ok {
		return domain.ActionPlan{}, fmt.Errorf("VM %q not found", request.VMID)
	}

	snapshot, ok := findSnapshot(vm.Snapshots, request.SnapshotID, request.SnapshotName)
	if !ok {
		return domain.ActionPlan{}, fmt.Errorf("snapshot %q not found for VM %q", snapshotRef(request), vm.ID)
	}
	operator := strings.TrimSpace(request.Operator)
	if operator == "" {
		return domain.ActionPlan{}, errors.New("operator is required")
	}

	plan := domain.ActionPlan{
		Action:            domain.ActionDeleteSnapshot,
		ConfirmationLevel: domain.ConfirmationTypedTargetName,
		ConfirmationToken: typedSnapshotDeleteToken(snapshot.Name),
		Operator:          operator,
		TargetHost:        vm.HostName,
		TargetType:        targetTypeVM,
		TargetID:          vm.ID,
		TargetName:        vm.Name,
		SnapshotID:        snapshot.ID,
		SnapshotName:      snapshot.Name,
		DeleteChildren:    request.DeleteChildren,
		Warnings:          []string{"Deleting a snapshot can consolidate disks and is not instantly reversible."},
	}
	plan.ID = planID(plan, plan.SnapshotID, plan.SnapshotName)

	return plan, nil
}

func (p Planner) inventory(ctx context.Context) (domain.Inventory, error) {
	if p.Provider == nil {
		return domain.Inventory{}, errors.New("snapshot planner requires a provider")
	}

	return p.Provider.Inventory(ctx)
}

func findVM(vms []domain.VM, id string) (domain.VM, bool) {
	for _, vm := range vms {
		if vm.ID == id {
			return vm, true
		}
	}

	return domain.VM{}, false
}

func findSnapshot(snapshots []domain.Snapshot, id, name string) (domain.Snapshot, bool) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	for _, snapshot := range snapshots {
		if id != "" && snapshot.ID == id {
			return snapshot, true
		}
		if name != "" && snapshot.Name == name {
			return snapshot, true
		}
	}

	return domain.Snapshot{}, false
}

func snapshotRef(request SnapshotDeleteRequest) string {
	if strings.TrimSpace(request.SnapshotID) != "" {
		return request.SnapshotID
	}
	return request.SnapshotName
}

func planID(plan domain.ActionPlan, parts ...string) string {
	material := []string{
		string(plan.Action),
		plan.TargetHost,
		plan.TargetType,
		plan.TargetID,
		plan.TargetName,
	}
	material = append(material, parts...)

	sum := sha256.Sum256([]byte(strings.Join(material, "\x00")))
	return "plan-" + hex.EncodeToString(sum[:])[:16]
}

func typedSnapshotDeleteToken(snapshotName string) string {
	return fmt.Sprintf("%s %s", strings.ToUpper(string(domain.ActionDeleteSnapshot)), strings.TrimSpace(snapshotName))
}
