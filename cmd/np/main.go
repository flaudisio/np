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
		Short: "CLI for deploying Nomad Pack applications from deploy.yaml",
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "deploy.yaml", "Path to deploy.yaml")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Print commands without executing")

	rootCmd.AddCommand(deployCmd())
	rootCmd.AddCommand(planCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func deployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a Nomad Pack from deploy.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("run")
		},
	}
}

func planCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Plan a Nomad Pack deployment from deploy.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("plan")
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
