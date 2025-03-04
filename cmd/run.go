package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/loft-sh/devpod-provider-kubernetes/pkg/kubernetes"
	"github.com/loft-sh/devpod-provider-kubernetes/pkg/options"
	"github.com/loft-sh/devpod/pkg/driver"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// RunCmd holds the cmd flags
type RunCmd struct{}

// NewRunCmd defines a command
func NewRunCmd() *cobra.Command {
	cmd := &RunCmd{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a container",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv()
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return runCmd
}

// Run runs the command logic
func (cmd *RunCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	var runOptions *driver.RunOptions

	runOptsEnv := os.Getenv("DEVCONTAINER_RUN_OPTIONS")
	if runOptsEnv != "" && runOptsEnv != "null" {
		runOptions = &driver.RunOptions{}
		err := json.Unmarshal([]byte(runOptsEnv), runOptions)
		if err != nil {
			return fmt.Errorf("unmarshal run options: %w", err)
		}
	}

	driver, err := kubernetes.NewKubernetesDriver(options, log)
	if err != nil {
		return err
	}

	return driver.RunDevContainer(ctx, options.DevContainerID, runOptions)
}
