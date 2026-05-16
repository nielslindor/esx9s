package provider

import (
	"context"

	"github.com/nielslindor/esx9s/internal/domain"
)

// Provider exposes a read-only inventory surface for standalone ESXi backends.
type Provider interface {
	Inventory(ctx context.Context) (domain.Inventory, error)
	Hosts(ctx context.Context) ([]domain.Host, error)
	VMs(ctx context.Context) ([]domain.VM, error)
	Datastores(ctx context.Context) ([]domain.Datastore, error)
	Tasks(ctx context.Context) ([]domain.Task, error)
}
