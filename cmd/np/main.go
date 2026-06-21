package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/flaudisio/np/internal/config"
	"github.com/flaudisio/np/internal/log"
	"github.com/flaudisio/np/internal/nomadpack"
)

var (
	configPath string
	dryRun     bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "np",
		Short: "CLI for deploying Nomad Pack applications from deploy.yml",
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "deploy.yml", "Path to deploy.yml")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "Print commands without executing")

	rootCmd.AddCommand(deployCmd())
	rootCmd.AddCommand(planCmd())
	rootCmd.AddCommand(destroyCmd())
	rootCmd.AddCommand(stopCmd())
	rootCmd.AddCommand(renderCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func deployCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"run"},
		Short:   "Deploy a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("run")
		},
	}
}

func planCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Plan a Nomad Pack deployment from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("plan")
		},
	}
}

func destroyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("destroy")
		},
	}
}

func stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("stop")
		},
	}
}

func renderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "render",
		Short: "Render a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("render")
		},
	}
}

func run(action string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	log.Info("Running from " + cwd)

	if err := nomadpack.Run(cfg, action, dryRun); err != nil {
		return err
	}

	log.Success(action + " complete")
	return nil
}
