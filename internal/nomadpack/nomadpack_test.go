package nomadpack

import (
	"os/exec"
	"testing"

	"github.com/flaudisio/np/internal/config"
)

func TestBuildCommandMinimal(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	cmd := BuildCommand(cfg, "run")
	expected := []string{"nomad-pack", "run", "my-pack"}
	assertSliceEqual(t, expected, cmd)
}

func TestBuildCommandWithDeployName(t *testing.T) {
	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{
			Name: "my-deployment",
		},
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	cmd := BuildCommand(cfg, "plan")
	expected := []string{"nomad-pack", "plan", "my-pack", "--name", "my-deployment"}
	assertSliceEqual(t, expected, cmd)
}

func TestBuildCommandWithVars(t *testing.T) {
	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{
			Vars: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	cmd := BuildCommand(cfg, "run")
	expected := []string{
		"nomad-pack", "run", "my-pack",
		"--var", "key1=value1",
		"--var", "key2=value2",
	}
	assertSliceEqual(t, expected, cmd)
}

func TestBuildCommandWithVarFiles(t *testing.T) {
	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{
			VarFiles: []string{"vars/a.yaml", "vars/b.yaml"},
		},
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	cmd := BuildCommand(cfg, "run")
	expected := []string{
		"nomad-pack", "run", "my-pack",
		"--var-file", "vars/a.yaml",
		"--var-file", "vars/b.yaml",
	}
	assertSliceEqual(t, expected, cmd)
}

func TestBuildCommandWithRegistry(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "my-registry",
				Ref:  "v1.0",
			},
		},
	}

	cmd := BuildCommand(cfg, "run")
	expected := []string{
		"nomad-pack", "run", "my-pack",
		"--registry", "my-registry",
		"--ref", "v1.0",
	}
	assertSliceEqual(t, expected, cmd)
}

func TestBuildCommandWithRegistryNoRef(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "my-registry",
			},
		},
	}

	cmd := BuildCommand(cfg, "plan")
	expected := []string{
		"nomad-pack", "plan", "my-pack",
		"--registry", "my-registry",
	}
	assertSliceEqual(t, expected, cmd)
}

func TestBuildCommandPlanVerbose(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
		Plan: config.PlanConfig{
			Verbose: true,
		},
	}

	cmd := BuildCommand(cfg, "plan")
	if !containsArg(cmd, "--verbose") {
		t.Error("expected --verbose flag in plan command")
	}
}

func TestBuildCommandPlanNotVerbose(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
		Plan: config.PlanConfig{
			Verbose: false,
		},
	}

	cmd := BuildCommand(cfg, "plan")
	if containsArg(cmd, "--verbose") {
		t.Error("expected no --verbose flag when verbose is false")
	}
}

func TestEnsureRegistryAlreadyExists(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME  SOURCE\nmy-registry  https://example.com\nother-reg  https://other.com\n")
		}
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	reg := &config.RegistryConfig{Name: "my-registry", Source: "https://example.com", Ref: "main"}
	err := ensureRegistry(reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 1 || calls[0][2] != "list" {
		t.Errorf("expected only list call, got: %v", calls)
	}
}

func TestEnsureRegistryNotExists(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME  SOURCE\nother-reg  https://other.com\n")
		}
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	reg := &config.RegistryConfig{Name: "my-registry", Source: "https://example.com", Ref: "main"}
	err := ensureRegistry(reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) < 2 {
		t.Fatalf("expected at least 2 calls (list + add), got %d: %v", len(calls), calls)
	}
	if calls[0][2] != "list" {
		t.Errorf("expected first call to be list, got: %v", calls[0])
	}
	if calls[1][2] != "add" {
		t.Errorf("expected second call to be add, got: %v", calls[1])
	}
}

func TestEnsureRegistryNoFalsePositive(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME  SOURCE\ncorpsec  https://example.com\n")
		}
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	reg := &config.RegistryConfig{Name: "corp", Source: "https://example.com", Ref: "main"}
	err := ensureRegistry(reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) < 2 {
		t.Fatalf("expected list + add calls (corp should not match corpsec), got %d: %v", len(calls), calls)
	}
}

func TestRunWithoutDryRun(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME  SOURCE\nmy-registry  https://example.com\n")
		}
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name:   "my-registry",
				Source: "https://example.com",
				Ref:    "main",
			},
		},
	}

	err := Run(cfg, "run", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) < 2 {
		t.Fatalf("expected at least 2 calls (list + run), got %d: %v", len(calls), calls)
	}
}

func TestRunDryRun(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	callCount := 0
	execCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	err := Run(cfg, "run", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 0 {
		t.Errorf("expected 0 exec calls, got %d", callCount)
	}
}

func assertSliceEqual(t *testing.T, expected, actual []string) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Fatalf("length mismatch: expected %v, got %v", expected, actual)
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("mismatch at index %d: expected %q, got %q\nfull expected: %v\nfull actual:   %v",
				i, expected[i], actual[i], expected, actual)
		}
	}
}

func containsArg(args []string, arg string) bool {
	for _, a := range args {
		if a == arg {
			return true
		}
	}
	return false
}
