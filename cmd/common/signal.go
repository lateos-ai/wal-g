package common

import (
	"os/signal"
	"syscall"
	"os"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// ExecuteContext runs cmd with a context canceled on SIGINT/SIGTERM. Exits
// with status 1 on error.
func ExecuteContext(cmd *cobra.Command) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
