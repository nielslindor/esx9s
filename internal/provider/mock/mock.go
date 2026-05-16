package mock

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nielslindor/esx9s/internal/domain"
)

const gib = int64(1024 * 1024 * 1024)

// Provider returns deterministic in-memory inventory for local development and tests.
type Provider struct {
	mu        sync.Mutex
	inventory domain.Inventory
	now       func() time.Time
}

// New returns a mock provider with multiple standalone ESXi hosts and VMs.
func New() *Provider {
	return &Provider{
		inventory: sampleInventory(),
		now:       func() time.Time { return time.Now().UTC() },
	}
}

// WithInventory returns a mock provider backed by caller-supplied inventory.
func WithInventory(inventory domain.Inventory) *Provider {
	return &Provider{
		inventory: cloneInventory(inventory),
		now:       func() time.Time { return time.Now().UTC() },
	}
}

func (p *Provider) Inventory(ctx context.Context) (domain.Inventory, error) {
	if err := ctx.Err(); err != nil {
		return domain.Inventory{}, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	return cloneInventory(p.inventory), nil
}

func (p *Provider) Hosts(ctx context.Context) ([]domain.Host, error) {
	inventory, err := p.Inventory(ctx)
	if err != nil {
		return nil, err
	}

	return inventory.Hosts, nil
}

func (p *Provider) VMs(ctx context.Context) ([]domain.VM, error) {
	inventory, err := p.Inventory(ctx)
	if err != nil {
		return nil, err
	}

	return inventory.VMs, nil
}

func (p *Provider) Datastores(ctx context.Context) ([]domain.Datastore, error) {
	inventory, err := p.Inventory(ctx)
	if err != nil {
		return nil, err
	}

	return inventory.Datastores, nil
}

func (p *Provider) Tasks(ctx context.Context) ([]domain.Task, error) {
	inventory, err := p.Inventory(ctx)
	if err != nil {
		return nil, err
	}

	return inventory.Tasks, nil
}

// ApplySnapshotPlan simulates a confirmed snapshot create/delete action.
//
// This method intentionally performs no ESXi or govmomi calls. It only updates
// deterministic in-memory state so UI flows can exercise plan/confirm/apply/audit.
func (p *Provider) ApplySnapshotPlan(ctx context.Context, plan domain.ActionPlan, confirmation string) (domain.ActionResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.ActionResult{}, err
	}
	if plan.Action != domain.ActionCreateSnapshot && plan.Action != domain.ActionDeleteSnapshot {
		return domain.ActionResult{}, fmt.Errorf("unsupported mock snapshot action %q", plan.Action)
	}
	if !plan.ConfirmedBy(confirmation) {
		return domain.ActionResult{}, errors.New("snapshot action confirmation did not match plan")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	vmIndex := -1
	for i, vm := range p.inventory.VMs {
		if vm.ID == plan.TargetID {
			vmIndex = i
			break
		}
	}
	if vmIndex < 0 {
		return domain.ActionResult{}, fmt.Errorf("VM %q not found", plan.TargetID)
	}

	now := p.currentTime()
	switch plan.Action {
	case domain.ActionCreateSnapshot:
		snapshot := domain.Snapshot{
			ID:          mockSnapshotID(plan.TargetID, plan.SnapshotName),
			Name:        plan.SnapshotName,
			Description: plan.SnapshotDesc,
			CreatedAt:   now,
			Current:     true,
		}
		for i := range p.inventory.VMs[vmIndex].Snapshots {
			p.inventory.VMs[vmIndex].Snapshots[i].Current = false
		}
		p.inventory.VMs[vmIndex].Snapshots = append(p.inventory.VMs[vmIndex].Snapshots, snapshot)
		plan.SnapshotID = snapshot.ID
	case domain.ActionDeleteSnapshot:
		snapshots := p.inventory.VMs[vmIndex].Snapshots
		deleteIndex := -1
		for i, snapshot := range snapshots {
			if snapshot.ID == plan.SnapshotID {
				deleteIndex = i
				break
			}
		}
		if deleteIndex < 0 {
			return domain.ActionResult{}, fmt.Errorf("snapshot %q not found", plan.SnapshotID)
		}
		p.inventory.VMs[vmIndex].Snapshots = append(snapshots[:deleteIndex], snapshots[deleteIndex+1:]...)
	}

	task := domain.Task{
		ID:         fmt.Sprintf("task-mock-%s", strings.TrimPrefix(plan.ID, "plan-")),
		HostID:     p.inventory.VMs[vmIndex].HostID,
		HostName:   p.inventory.VMs[vmIndex].HostName,
		TargetID:   plan.TargetID,
		TargetName: plan.TargetName,
		Operation:  string(plan.Action),
		Status:     domain.TaskStatusSuccess,
		StartedAt:  now,
		FinishedAt: &now,
	}
	p.inventory.Tasks = append(p.inventory.Tasks, task)

	return domain.ActionResult{
		PlanID: plan.ID,
		Action: plan.Action,
		Task:   task,
		AuditEvent: domain.AuditEvent{
			Timestamp:  now,
			Operator:   plan.Operator,
			TargetHost: plan.TargetHost,
			TargetType: plan.TargetType,
			TargetID:   plan.TargetID,
			TargetName: plan.TargetName,
			Action:     string(plan.Action),
			PlanID:     plan.ID,
			Result:     "mock_success",
		},
	}, nil
}

func sampleInventory() domain.Inventory {
	finished := time.Date(2026, 5, 16, 9, 5, 0, 0, time.UTC)

	return domain.Inventory{
		Hosts: []domain.Host{
			{
				ID:              "host-esxi01",
				Name:            "esxi01",
				Address:         "192.168.10.21",
				Product:         "VMware ESXi",
				Version:         "8.0.2",
				Status:          domain.HostStatusConnected,
				CPUCapacityMHz:  38400,
				CPUUsedMHz:      9200,
				MemoryBytes:     128 * gib,
				MemoryUsedBytes: 72 * gib,
				VMCount:         3,
			},
			{
				ID:              "host-esxi02",
				Name:            "esxi02",
				Address:         "192.168.10.22",
				Product:         "VMware ESXi",
				Version:         "8.0.2",
				Status:          domain.HostStatusMaintenance,
				CPUCapacityMHz:  28800,
				CPUUsedMHz:      4100,
				MemoryBytes:     96 * gib,
				MemoryUsedBytes: 31 * gib,
				VMCount:         2,
			},
		},
		VMs: []domain.VM{
			{
				ID:           "vm-home-assistant",
				Name:         "home-assistant",
				HostID:       "host-esxi01",
				HostName:     "esxi01",
				DatastoreIDs: []string{"ds-esxi01-nvme"},
				Snapshots: []domain.Snapshot{
					{
						ID:          "snap-home-assistant-before-upgrade",
						Name:        "before-upgrade",
						Description: "before OS package upgrade",
						CreatedAt:   time.Date(2026, 5, 15, 20, 15, 0, 0, time.UTC),
						Current:     true,
					},
				},
				Path:             "[nvme-local-01] home-assistant/home-assistant.vmx",
				PowerState:       domain.PowerStatePoweredOn,
				GuestOS:          "Ubuntu Linux (64-bit)",
				IPAddress:        "192.168.30.10",
				CPUCount:         2,
				MemoryBytes:      4 * gib,
				ProvisionedBytes: 64 * gib,
				UsedBytes:        29 * gib,
				UptimeSeconds:    604800,
			},
			{
				ID:               "vm-gitlab-runner-01",
				Name:             "gitlab-runner-01",
				HostID:           "host-esxi01",
				HostName:         "esxi01",
				DatastoreIDs:     []string{"ds-esxi01-nvme"},
				Path:             "[nvme-local-01] gitlab-runner-01/gitlab-runner-01.vmx",
				PowerState:       domain.PowerStatePoweredOn,
				GuestOS:          "Debian GNU/Linux 12 (64-bit)",
				IPAddress:        "192.168.30.31",
				CPUCount:         4,
				MemoryBytes:      8 * gib,
				ProvisionedBytes: 120 * gib,
				UsedBytes:        47 * gib,
				UptimeSeconds:    172800,
			},
			{
				ID:               "vm-template-debian12",
				Name:             "template-debian12",
				HostID:           "host-esxi01",
				HostName:         "esxi01",
				DatastoreIDs:     []string{"ds-esxi01-archive"},
				Path:             "[archive-01] template-debian12/template-debian12.vmx",
				PowerState:       domain.PowerStatePoweredOff,
				GuestOS:          "Debian GNU/Linux 12 (64-bit)",
				CPUCount:         2,
				MemoryBytes:      2 * gib,
				ProvisionedBytes: 40 * gib,
				UsedBytes:        12 * gib,
			},
			{
				ID:               "vm-truenas",
				Name:             "truenas",
				HostID:           "host-esxi02",
				HostName:         "esxi02",
				DatastoreIDs:     []string{"ds-esxi02-ssd"},
				Path:             "[ssd-local-02] truenas/truenas.vmx",
				PowerState:       domain.PowerStatePoweredOn,
				GuestOS:          "FreeBSD 13 or later (64-bit)",
				IPAddress:        "192.168.30.20",
				CPUCount:         4,
				MemoryBytes:      16 * gib,
				ProvisionedBytes: 80 * gib,
				UsedBytes:        38 * gib,
				UptimeSeconds:    259200,
			},
			{
				ID:               "vm-windows-lab",
				Name:             "windows-lab",
				HostID:           "host-esxi02",
				HostName:         "esxi02",
				DatastoreIDs:     []string{"ds-esxi02-ssd"},
				Path:             "[ssd-local-02] windows-lab/windows-lab.vmx",
				PowerState:       domain.PowerStateSuspended,
				GuestOS:          "Microsoft Windows 11 (64-bit)",
				CPUCount:         4,
				MemoryBytes:      8 * gib,
				ProvisionedBytes: 100 * gib,
				UsedBytes:        61 * gib,
			},
		},
		Datastores: []domain.Datastore{
			{
				ID:            "ds-esxi01-nvme",
				Name:          "nvme-local-01",
				HostID:        "host-esxi01",
				HostName:      "esxi01",
				Type:          "VMFS",
				CapacityBytes: 1800 * gib,
				FreeBytes:     980 * gib,
				Accessible:    true,
			},
			{
				ID:            "ds-esxi01-archive",
				Name:          "archive-01",
				HostID:        "host-esxi01",
				HostName:      "esxi01",
				Type:          "NFS",
				CapacityBytes: 4096 * gib,
				FreeBytes:     2750 * gib,
				Accessible:    true,
			},
			{
				ID:            "ds-esxi02-ssd",
				Name:          "ssd-local-02",
				HostID:        "host-esxi02",
				HostName:      "esxi02",
				Type:          "VMFS",
				CapacityBytes: 960 * gib,
				FreeBytes:     340 * gib,
				Accessible:    true,
			},
		},
		Tasks: []domain.Task{
			{
				ID:         "task-1001",
				HostID:     "host-esxi01",
				HostName:   "esxi01",
				TargetID:   "vm-gitlab-runner-01",
				TargetName: "gitlab-runner-01",
				Operation:  "PowerOnVM",
				Status:     domain.TaskStatusSuccess,
				StartedAt:  time.Date(2026, 5, 16, 9, 4, 30, 0, time.UTC),
				FinishedAt: &finished,
			},
			{
				ID:         "task-1002",
				HostID:     "host-esxi02",
				HostName:   "esxi02",
				TargetID:   "host-esxi02",
				TargetName: "esxi02",
				Operation:  "EnterMaintenanceMode",
				Status:     domain.TaskStatusRunning,
				StartedAt:  time.Date(2026, 5, 16, 9, 7, 0, 0, time.UTC),
			},
		},
	}
}

func cloneInventory(inventory domain.Inventory) domain.Inventory {
	return domain.Inventory{
		Hosts:      append([]domain.Host(nil), inventory.Hosts...),
		VMs:        cloneVMs(inventory.VMs),
		Datastores: append([]domain.Datastore(nil), inventory.Datastores...),
		Tasks:      append([]domain.Task(nil), inventory.Tasks...),
	}
}

func cloneVMs(vms []domain.VM) []domain.VM {
	cloned := make([]domain.VM, len(vms))
	for i, vm := range vms {
		cloned[i] = vm
		cloned[i].DatastoreIDs = append([]string(nil), vm.DatastoreIDs...)
		cloned[i].Snapshots = append([]domain.Snapshot(nil), vm.Snapshots...)
	}

	return cloned
}

func mockSnapshotID(vmID, name string) string {
	id := strings.ToLower(strings.TrimSpace(vmID + "-" + name))
	id = strings.NewReplacer(" ", "-", "_", "-", "/", "-", "\\", "-").Replace(id)
	return "snap-" + id
}

func (p *Provider) currentTime() time.Time {
	if p.now == nil {
		return time.Now().UTC()
	}
	return p.now().UTC()
}
