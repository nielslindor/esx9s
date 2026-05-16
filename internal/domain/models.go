package domain

import (
	"strings"
	"time"
)

// HostStatus is the high-level connection and health state for an ESXi host.
type HostStatus string

const (
	HostStatusUnknown      HostStatus = "unknown"
	HostStatusConnected    HostStatus = "connected"
	HostStatusDisconnected HostStatus = "disconnected"
	HostStatusMaintenance  HostStatus = "maintenance"
)

// PowerState is the runtime power state for a virtual machine.
type PowerState string

const (
	PowerStateUnknown    PowerState = "unknown"
	PowerStatePoweredOn  PowerState = "powered_on"
	PowerStateSuspended  PowerState = "suspended"
	PowerStatePoweredOff PowerState = "powered_off"
)

// TaskStatus is the lifecycle state for a provider task.
type TaskStatus string

const (
	TaskStatusQueued  TaskStatus = "queued"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusSuccess TaskStatus = "success"
	TaskStatusError   TaskStatus = "error"
)

// Inventory is a point-in-time view of all provider-visible resources.
type Inventory struct {
	Hosts      []Host      `json:"hosts"`
	VMs        []VM        `json:"vms"`
	Datastores []Datastore `json:"datastores"`
	Tasks      []Task      `json:"tasks"`
}

// Host describes one standalone ESXi host managed by esx9s.
type Host struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Address         string     `json:"address"`
	Product         string     `json:"product"`
	Version         string     `json:"version"`
	Status          HostStatus `json:"status"`
	CPUCapacityMHz  int64      `json:"cpu_capacity_mhz"`
	CPUUsedMHz      int64      `json:"cpu_used_mhz"`
	MemoryBytes     int64      `json:"memory_bytes"`
	MemoryUsedBytes int64      `json:"memory_used_bytes"`
	VMCount         int        `json:"vm_count"`
}

// VM describes a virtual machine in the unified inventory.
type VM struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	HostID           string     `json:"host_id"`
	HostName         string     `json:"host_name"`
	DatastoreIDs     []string   `json:"datastore_ids"`
	Snapshots        []Snapshot `json:"snapshots,omitempty"`
	Path             string     `json:"path"`
	PowerState       PowerState `json:"power_state"`
	GuestOS          string     `json:"guest_os"`
	IPAddress        string     `json:"ip_address,omitempty"`
	CPUCount         int        `json:"cpu_count"`
	MemoryBytes      int64      `json:"memory_bytes"`
	ProvisionedBytes int64      `json:"provisioned_bytes"`
	UsedBytes        int64      `json:"used_bytes"`
	UptimeSeconds    int64      `json:"uptime_seconds,omitempty"`
}

// Snapshot describes VM snapshot metadata visible in inventory and action plans.
type Snapshot struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	Current     bool      `json:"current,omitempty"`
}

// Datastore describes storage visible to a standalone ESXi host.
type Datastore struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	HostID        string `json:"host_id"`
	HostName      string `json:"host_name"`
	Type          string `json:"type"`
	CapacityBytes int64  `json:"capacity_bytes"`
	FreeBytes     int64  `json:"free_bytes"`
	Accessible    bool   `json:"accessible"`
}

// Task describes recent provider-side activity.
type Task struct {
	ID         string     `json:"id"`
	HostID     string     `json:"host_id"`
	HostName   string     `json:"host_name"`
	TargetID   string     `json:"target_id,omitempty"`
	TargetName string     `json:"target_name,omitempty"`
	Operation  string     `json:"operation"`
	Status     TaskStatus `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Error      string     `json:"error,omitempty"`
}

// ActionType is an auditable operator action that can be planned and applied.
type ActionType string

const (
	ActionCreateSnapshot ActionType = "snapshot_create"
	ActionDeleteSnapshot ActionType = "snapshot_delete"
)

// ConfirmationLevel describes the operator confirmation required before apply.
type ConfirmationLevel string

const (
	ConfirmationVisibleTarget   ConfirmationLevel = "visible_target"
	ConfirmationTypedTargetName ConfirmationLevel = "typed_target_name"
)

// ActionPlan is the immutable, audit-friendly description shown before apply.
type ActionPlan struct {
	ID                string            `json:"id"`
	Action            ActionType        `json:"action"`
	ConfirmationLevel ConfirmationLevel `json:"confirmation_level"`
	ConfirmationToken string            `json:"confirmation_token,omitempty"`
	Operator          string            `json:"operator"`
	TargetHost        string            `json:"target_host"`
	TargetType        string            `json:"target_type"`
	TargetID          string            `json:"target_id"`
	TargetName        string            `json:"target_name"`
	SnapshotID        string            `json:"snapshot_id,omitempty"`
	SnapshotName      string            `json:"snapshot_name,omitempty"`
	SnapshotDesc      string            `json:"snapshot_description,omitempty"`
	IncludeMemory     bool              `json:"include_memory,omitempty"`
	Quiesce           bool              `json:"quiesce,omitempty"`
	DeleteChildren    bool              `json:"delete_children,omitempty"`
	Warnings          []string          `json:"warnings,omitempty"`
}

// ConfirmedBy reports whether input satisfies the confirmation required by the plan.
func (p ActionPlan) ConfirmedBy(input string) bool {
	switch p.ConfirmationLevel {
	case ConfirmationVisibleTarget:
		return strings.EqualFold(strings.TrimSpace(input), "yes")
	case ConfirmationTypedTargetName:
		return strings.TrimSpace(input) == p.ConfirmationToken && p.ConfirmationToken != ""
	default:
		return false
	}
}

// ActionResult is returned after an action executor applies a confirmed plan.
type ActionResult struct {
	PlanID     string     `json:"plan_id"`
	Action     ActionType `json:"action"`
	Task       Task       `json:"task"`
	AuditEvent AuditEvent `json:"audit_event"`
}

// AuditEvent mirrors the minimum append-only fields without coupling domain to storage.
type AuditEvent struct {
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
