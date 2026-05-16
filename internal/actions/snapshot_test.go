package actions

import (
	"context"
	"strings"
	"testing"

	"github.com/nielslindor/esx9s/internal/domain"
	"github.com/nielslindor/esx9s/internal/provider/mock"
)

func TestPlanSnapshotCreateBuildsVisibleConfirmationPlan(t *testing.T) {
	planner := Planner{Provider: mock.New()}

	plan, err := planner.PlanSnapshotCreate(context.Background(), SnapshotCreateRequest{
		Operator:      "local-user",
		VMID:          "vm-home-assistant",
		Name:          "before-upgrade",
		Description:   "before OS package upgrade",
		IncludeMemory: true,
	})
	if err != nil {
		t.Fatalf("PlanSnapshotCreate() error = %v", err)
	}

	if plan.Action != domain.ActionCreateSnapshot {
		t.Fatalf("action = %q, want %q", plan.Action, domain.ActionCreateSnapshot)
	}
	if plan.ConfirmationLevel != domain.ConfirmationVisibleTarget {
		t.Fatalf("confirmation level = %q, want %q", plan.ConfirmationLevel, domain.ConfirmationVisibleTarget)
	}
	if !plan.ConfirmedBy("yes") {
		t.Fatal("create plan was not confirmed by visible target confirmation")
	}
	if plan.TargetName != "home-assistant" || plan.TargetHost != "esxi01" {
		t.Fatalf("target = %s/%s, want esxi01/home-assistant", plan.TargetHost, plan.TargetName)
	}
	if plan.SnapshotName != "before-upgrade" {
		t.Fatalf("snapshot name = %q, want before-upgrade", plan.SnapshotName)
	}
	if plan.ID == "" {
		t.Fatal("plan ID was empty")
	}
}

func TestPlanSnapshotDeleteRequiresTypedSnapshotName(t *testing.T) {
	planner := Planner{Provider: mock.New()}

	plan, err := planner.PlanSnapshotDelete(context.Background(), SnapshotDeleteRequest{
		Operator:   "local-user",
		VMID:       "vm-home-assistant",
		SnapshotID: "snap-home-assistant-before-upgrade",
	})
	if err != nil {
		t.Fatalf("PlanSnapshotDelete() error = %v", err)
	}

	if plan.Action != domain.ActionDeleteSnapshot {
		t.Fatalf("action = %q, want %q", plan.Action, domain.ActionDeleteSnapshot)
	}
	if plan.ConfirmationLevel != domain.ConfirmationTypedTargetName {
		t.Fatalf("confirmation level = %q, want %q", plan.ConfirmationLevel, domain.ConfirmationTypedTargetName)
	}
	if plan.ConfirmedBy("confirm") {
		t.Fatal("delete plan accepted generic confirmation")
	}
	if !plan.ConfirmedBy("SNAPSHOT_DELETE before-upgrade") {
		t.Fatal("delete plan did not accept exact snapshot name")
	}
	if len(plan.Warnings) == 0 {
		t.Fatal("delete plan did not include operator warning")
	}
}

func TestPlanSnapshotCreateRejectsMissingSnapshotName(t *testing.T) {
	planner := Planner{Provider: mock.New()}

	_, err := planner.PlanSnapshotCreate(context.Background(), SnapshotCreateRequest{
		Operator: "local-user",
		VMID:     "vm-home-assistant",
		Name:     " ",
	})
	if err == nil {
		t.Fatal("PlanSnapshotCreate() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "snapshot name") {
		t.Fatalf("error = %q, want snapshot name validation", err)
	}
}
