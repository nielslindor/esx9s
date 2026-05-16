package mock

import (
	"context"
	"errors"
	"testing"

	"github.com/nielslindor/esx9s/internal/actions"
	"github.com/nielslindor/esx9s/internal/domain"
	"github.com/nielslindor/esx9s/internal/provider"
)

func TestProviderImplementsProviderInterface(t *testing.T) {
	var _ provider.Provider = New()
}

func TestInventoryContainsHostsVMsDatastoresAndTasks(t *testing.T) {
	inventory, err := New().Inventory(context.Background())
	if err != nil {
		t.Fatalf("Inventory returned error: %v", err)
	}

	if got, want := len(inventory.Hosts), 2; got != want {
		t.Fatalf("host count = %d, want %d", got, want)
	}
	if got, want := len(inventory.VMs), 5; got != want {
		t.Fatalf("vm count = %d, want %d", got, want)
	}
	if got, want := len(inventory.Datastores), 3; got != want {
		t.Fatalf("datastore count = %d, want %d", got, want)
	}
	if got, want := len(inventory.Tasks), 2; got != want {
		t.Fatalf("task count = %d, want %d", got, want)
	}

	hostIDs := map[string]bool{}
	for _, host := range inventory.Hosts {
		hostIDs[host.ID] = true
	}

	for _, vm := range inventory.VMs {
		if !hostIDs[vm.HostID] {
			t.Fatalf("VM %q references unknown host %q", vm.Name, vm.HostID)
		}
		if len(vm.DatastoreIDs) == 0 {
			t.Fatalf("VM %q has no datastores", vm.Name)
		}
	}
}

func TestInventoryReturnsDefensiveCopies(t *testing.T) {
	provider := New()

	first, err := provider.Inventory(context.Background())
	if err != nil {
		t.Fatalf("first Inventory returned error: %v", err)
	}
	first.Hosts[0].Name = "changed"
	first.VMs[0].DatastoreIDs[0] = "changed"

	second, err := provider.Inventory(context.Background())
	if err != nil {
		t.Fatalf("second Inventory returned error: %v", err)
	}

	if second.Hosts[0].Name == "changed" {
		t.Fatal("Inventory returned shared host slice")
	}
	if second.VMs[0].DatastoreIDs[0] == "changed" {
		t.Fatal("Inventory returned shared VM datastore IDs")
	}
}

func TestInventoryHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := New().Inventory(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Inventory error = %v, want context.Canceled", err)
	}
}

func TestWithInventoryUsesCallerSuppliedData(t *testing.T) {
	mock := WithInventory(domain.Inventory{
		Hosts: []domain.Host{{ID: "host-test", Name: "test", Status: domain.HostStatusConnected}},
	})

	hosts, err := mock.Hosts(context.Background())
	if err != nil {
		t.Fatalf("Hosts returned error: %v", err)
	}

	if got, want := len(hosts), 1; got != want {
		t.Fatalf("host count = %d, want %d", got, want)
	}
	if got, want := hosts[0].ID, "host-test"; got != want {
		t.Fatalf("host ID = %q, want %q", got, want)
	}
}

func TestApplySnapshotCreatePlanRequiresConfirmationAndAudits(t *testing.T) {
	mock := New()
	planner := actions.Planner{Provider: mock}
	plan, err := planner.PlanSnapshotCreate(context.Background(), actions.SnapshotCreateRequest{
		Operator: "local-user",
		VMID:     "vm-template-debian12",
		Name:     "before-test",
	})
	if err != nil {
		t.Fatalf("PlanSnapshotCreate() error = %v", err)
	}

	if _, err := mock.ApplySnapshotPlan(context.Background(), plan, "wrong"); err == nil {
		t.Fatal("ApplySnapshotPlan() error = nil, want confirmation error")
	}

	result, err := mock.ApplySnapshotPlan(context.Background(), plan, "yes")
	if err != nil {
		t.Fatalf("ApplySnapshotPlan() error = %v", err)
	}

	if result.AuditEvent.PlanID != plan.ID {
		t.Fatalf("audit plan ID = %q, want %q", result.AuditEvent.PlanID, plan.ID)
	}
	if result.AuditEvent.Result != "mock_success" {
		t.Fatalf("audit result = %q, want mock_success", result.AuditEvent.Result)
	}

	inventory, err := mock.Inventory(context.Background())
	if err != nil {
		t.Fatalf("Inventory() error = %v", err)
	}

	vm := findTestVM(t, inventory.VMs, "vm-template-debian12")
	if got, want := len(vm.Snapshots), 1; got != want {
		t.Fatalf("snapshot count = %d, want %d", got, want)
	}
	if vm.Snapshots[0].Name != "before-test" {
		t.Fatalf("snapshot name = %q, want before-test", vm.Snapshots[0].Name)
	}
	if got, want := len(inventory.Tasks), 3; got != want {
		t.Fatalf("task count = %d, want %d", got, want)
	}
}

func TestApplySnapshotDeletePlanRemovesSnapshotAfterTypedConfirmation(t *testing.T) {
	mock := New()
	planner := actions.Planner{Provider: mock}
	plan, err := planner.PlanSnapshotDelete(context.Background(), actions.SnapshotDeleteRequest{
		Operator:   "local-user",
		VMID:       "vm-home-assistant",
		SnapshotID: "snap-home-assistant-before-upgrade",
	})
	if err != nil {
		t.Fatalf("PlanSnapshotDelete() error = %v", err)
	}

	result, err := mock.ApplySnapshotPlan(context.Background(), plan, "SNAPSHOT_DELETE before-upgrade")
	if err != nil {
		t.Fatalf("ApplySnapshotPlan() error = %v", err)
	}
	if result.Action != domain.ActionDeleteSnapshot {
		t.Fatalf("result action = %q, want %q", result.Action, domain.ActionDeleteSnapshot)
	}

	inventory, err := mock.Inventory(context.Background())
	if err != nil {
		t.Fatalf("Inventory() error = %v", err)
	}

	vm := findTestVM(t, inventory.VMs, "vm-home-assistant")
	if len(vm.Snapshots) != 0 {
		t.Fatalf("snapshot count = %d, want 0", len(vm.Snapshots))
	}
}

func findTestVM(t *testing.T, vms []domain.VM, id string) domain.VM {
	t.Helper()
	for _, vm := range vms {
		if vm.ID == id {
			return vm
		}
	}
	t.Fatalf("VM %q not found", id)
	return domain.VM{}
}
