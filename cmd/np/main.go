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
	cdDir      string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "np",
		Short: "CLI for deploying Nomad Pack applications from deploy.yml",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cdDir != "" {
				if err := os.Chdir(cdDir); err != nil {
					return fmt.Errorf("changing directory to %s: %w", cdDir, err)
				}
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "deploy.yml", "Path to deploy.yml")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "Print commands without executing")
	rootCmd.PersistentFlags().StringVarP(&cdDir, "cd", "C", "", "Change to directory before running commands")

	rootCmd.AddCommand(deployCmd())
	rootCmd.AddCommand(planCmd())
	rootCmd.AddCommand(destroyCmd())
	rootCmd.AddCommand(stopCmd())
	rootCmd.AddCommand(renderCmd())
	rootCmd.AddCommand(registryCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// deployCmd returns the "deploy" (aliased as "run") subcommand.
func deployCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"run"},
		Short:   "Deploy a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("run", args)
		},
	}
}

// planCmd returns the "plan" subcommand.
func planCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Plan a Nomad Pack deployment from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("plan", args)
		},
	}
}

// destroyCmd returns the "destroy" subcommand.
func destroyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("destroy", args)
		},
	}
}

// stopCmd returns the "stop" subcommand.
func stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("stop", args)
		},
	}
}

// renderCmd returns the "render" subcommand.
func renderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "render",
		Short: "Render a Nomad Pack from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("render", args)
		},
	}
}

// registryCmd returns the "registry" subcommand group.
func registryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Manage Nomad Pack registries from deploy.yml",
	}
	cmd.AddCommand(registryAddCmd())
	cmd.AddCommand(registryDeleteCmd())
	cmd.AddCommand(registryUpdateCmd())
	return cmd
}

// registryAddCmd returns the "add" subcommand.
func registryAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add the configured registry from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			return nomadpack.RegistryAdd(cfg, dryRun)
		},
	}
}

// registryDeleteCmd returns the "delete" subcommand.
func registryDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete the configured registry from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			return nomadpack.RegistryDelete(cfg, dryRun)
		},
	}
}

// registryUpdateCmd returns the "update" subcommand.
func registryUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the configured registry from deploy.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			return nomadpack.RegistryUpdate(cfg, dryRun)
		},
	}
}

// run loads the config, then builds and executes the nomad-pack action.
func run(action string, extraArgs []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	log.Info("Running from " + cwd)

	if err := nomadpack.Run(cfg, action, dryRun, extraArgs...); err != nil {
		return err
	}

	log.Info(action + " complete")
	return nil
}
