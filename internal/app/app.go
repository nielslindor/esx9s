package app

import (
	"context"
	"io"

	"github.com/nielslindor/esx9s/internal/domain"
	"github.com/nielslindor/esx9s/internal/tui"
)

// Run starts the interactive operator console.
func Run(ctx context.Context, stdin io.Reader, stdout io.Writer) error {
	return tui.Run(ctx, stdin, stdout)
}

// RunWithInventory starts the operator console with a preloaded inventory.
func RunWithInventory(ctx context.Context, stdin io.Reader, stdout io.Writer, inventory domain.Inventory) error {
	return tui.RunModel(ctx, stdin, stdout, tui.NewModelFromInventory(inventory))
}

// RenderInventory writes the initial terminal view for non-interactive launches.
func RenderInventory(stdout io.Writer, inventory domain.Inventory) error {
	_, err := io.WriteString(stdout, tui.NewModelFromInventory(inventory).View())
	return err
}
