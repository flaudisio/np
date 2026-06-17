package nomadpack

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/flaudisio/np/internal/config"
	"github.com/flaudisio/np/internal/log"
)

var execCommand = exec.Command

func BuildCommand(cfg *config.DeployConfig, action string) []string {
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

	return cmd
}

func Run(cfg *config.DeployConfig, action string, dryRun bool) error {
	cmd := BuildCommand(cfg, action)
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

	return c.Run()
}

func SetExecCommand(fn func(name string, args ...string) *exec.Cmd) func(name string, args ...string) *exec.Cmd {
	old := execCommand
	execCommand = fn
	return old
}

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
