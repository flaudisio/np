// Package nomadpack builds and executes nomad-pack CLI commands.
package nomadpack

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/flaudisio/np/internal/config"
	"github.com/flaudisio/np/internal/log"
)

var execCommand = exec.Command

// BuildCommand builds a nomad-pack argument slice from [config.DeployConfig].
// Returns the full command line as []string suitable for exec.Command.
//
// Arguments are appended in this order: action, pack name, name, sorted vars,
// sorted var files, optional registry flags, --verbose (plan only), extraArgs.
func BuildCommand(cfg *config.DeployConfig, action string, extraArgs ...string) []string {
	cmd := []string{"nomad-pack", action, cfg.Pack.Name}

	if cfg.Deploy.Name != "" {
		cmd = append(cmd, "--name", cfg.Deploy.Name)
	}

	keys := make([]string, 0, len(cfg.Deploy.Vars))
	for k := range cfg.Deploy.Vars {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		cmd = append(cmd, "--var", fmt.Sprintf("%s=%s", k, cfg.Deploy.Vars[k]))
	}

	for _, f := range cfg.Deploy.VarFiles {
		cmd = append(cmd, "--var-file", f)
	}

	if cfg.Pack.Registry != nil {
		cmd = append(cmd, "--registry", cfg.Pack.Registry.Name)
		if cfg.Pack.Registry.Ref != "" {
			cmd = append(cmd, "--ref", cfg.Pack.Registry.Ref)
		}
	}

	if action == "plan" && cfg.Plan.Verbose {
		cmd = append(cmd, "--verbose")
	}

	cmd = append(cmd, extraArgs...)

	return cmd
}

// Run builds and executes the nomad-pack command.
// If dryRun is true, only the command is logged and nil is returned.
// Before running, the registry (if configured) is registered via ensureRegistry.
//
// Plan exit code 1 (non-empty plan) and run exit code 2 are treated as
// success — returned as nil error.
func Run(cfg *config.DeployConfig, action string, dryRun bool, extraArgs ...string) error {
	cmd := BuildCommand(cfg, action, extraArgs...)
	log.Info(fmt.Sprintf("+ %s", strings.Join(cmd, " ")))

	if dryRun {
		return nil
	}

	if cfg.Pack.Registry != nil {
		if err := ensureRegistry(cfg.Pack.Registry); err != nil {
			return fmt.Errorf("registry: %w", err)
		}
	}

	c := execCommand(cmd[0], cmd[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if action == "plan" && exitErr.ExitCode() == 1 {
				return nil
			}
			if action == "run" && exitErr.ExitCode() == 2 {
				return nil
			}
		}
		return err
	}
	return nil
}

// SetExecCommand replaces the function used to create exec.Cmd instances.
// Intended for tests. Returns the previous function for restoration.
func SetExecCommand(fn func(name string, args ...string) *exec.Cmd) func(name string, args ...string) *exec.Cmd {
	old := execCommand
	execCommand = fn
	return old
}

// ensureRegistry lists existing registries and adds the configured one
// if it is not already present.
func ensureRegistry(reg *config.RegistryConfig) error {
	out, err := execCommand("nomad-pack", "registry", "list").Output()
	if err != nil {
		return fmt.Errorf("listing registries: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == reg.Name {
			return nil
		}
	}

	log.Info(fmt.Sprintf("registering registry %s from %s", reg.Name, reg.Source))

	args := []string{"registry", "add", reg.Name, reg.Source}
	if reg.Ref != "" {
		args = append(args, "--ref", reg.Ref)
	}

	c := execCommand("nomad-pack", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
