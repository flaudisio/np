package nomadpack

import (
	"os/exec"
	"strings"
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

func TestBuildCommandWithExtraArgs(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	cmd := BuildCommand(cfg, "render", "--no-format", "--arg", "value")
	expected := []string{
		"nomad-pack", "render", "my-pack",
		"--no-format", "--arg", "value",
	}
	assertSliceEqual(t, expected, cmd)
}

func TestEnsureRegistryAlreadyExists(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME | REF | LOCAL REF | REGISTRY URL\nmy-registry | main | abc123 | https://example.com\nother-reg | latest | def456 | https://other.com\n")
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

func TestEnsureRegistryDifferentRef(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME | REF | LOCAL REF | REGISTRY URL\nmy-registry | old | abc123 | https://example.com\n")
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
		t.Fatalf("expected list + add calls (same name, new ref), got %d: %v", len(calls), calls)
	}
	if calls[1][2] != "add" || calls[1][5] != "--ref" || calls[1][6] != "main" {
		t.Errorf("expected add with --ref main, got: %v", calls[1])
	}
}

func TestEnsureRegistryDefaultRef(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME | REF | LOCAL REF | REGISTRY URL\nmy-registry | latest | abc123 | https://example.com\n")
		}
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	reg := &config.RegistryConfig{Name: "my-registry", Source: "https://example.com"}
	err := ensureRegistry(reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 1 || calls[0][2] != "list" {
		t.Errorf("expected only list call (empty ref matches latest), got: %v", calls)
	}
}

func TestEnsureRegistryNotExists(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "nomad-pack" && len(args) >= 2 && args[0] == "registry" && args[1] == "list" {
			return exec.Command("echo", "REGISTRY NAME | REF | LOCAL REF | REGISTRY URL\nother-reg | latest | def456 | https://other.com\n")
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
			return exec.Command("echo", "REGISTRY NAME | REF | LOCAL REF | REGISTRY URL\ncorpsec | latest | abc123 | https://example.com\n")
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
			return exec.Command("echo", "REGISTRY NAME | REF | LOCAL REF | REGISTRY URL\nmy-registry | main | abc123 | https://example.com\n")
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

func TestRegistryAdd(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name:   "community",
				Source: "https://example.com/registry",
				Ref:    "v0.1.0",
			},
		},
	}

	err := RegistryAdd(cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad-pack", "registry", "add", "community", "https://example.com/registry", "--ref", "v0.1.0"}
	assertSliceEqual(t, expected, calls[0])
}

func TestRegistryAddNoRef(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name:   "community",
				Source: "https://example.com/registry",
			},
		},
	}

	err := RegistryAdd(cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad-pack", "registry", "add", "community", "https://example.com/registry"}
	assertSliceEqual(t, expected, calls[0])
}

func TestRegistryAddDryRun(t *testing.T) {
	callCount := 0
	execCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name:   "community",
				Source: "https://example.com/registry",
			},
		},
	}

	err := RegistryAdd(cfg, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 0 {
		t.Errorf("expected 0 exec calls, got %d", callCount)
	}
}

func TestRegistryAddNoRegistry(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	err := RegistryAdd(cfg, false)
	if err == nil {
		t.Fatal("expected error for missing registry config")
	}
	if err.Error() != "no registry configured in deploy.yml" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegistryDelete(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
			},
		},
	}

	err := RegistryDelete(cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad-pack", "registry", "delete", "community"}
	assertSliceEqual(t, expected, calls[0])
}

func TestRegistryDeleteWithRef(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
				Ref:  "v0.1.0",
			},
		},
	}

	err := RegistryDelete(cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad-pack", "registry", "delete", "community", "--ref", "v0.1.0"}
	assertSliceEqual(t, expected, calls[0])
}

func TestRegistryDeleteDryRun(t *testing.T) {
	callCount := 0
	execCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
			},
		},
	}

	err := RegistryDelete(cfg, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 0 {
		t.Errorf("expected 0 exec calls, got %d", callCount)
	}
}

func TestRegistryDeleteNoRegistry(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	err := RegistryDelete(cfg, false)
	if err == nil {
		t.Fatal("expected error for missing registry config")
	}
	if err.Error() != "no registry configured in deploy.yml" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegistryUpdate(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
				Ref:  "v0.2.0",
			},
		},
	}

	err := RegistryUpdate(cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad-pack", "registry", "update", "community", "--ref", "v0.2.0"}
	assertSliceEqual(t, expected, calls[0])
}

func TestRegistryUpdateNoRef(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
			},
		},
	}

	err := RegistryUpdate(cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad-pack", "registry", "update", "community"}
	assertSliceEqual(t, expected, calls[0])
}

func TestRegistryUpdateDryRun(t *testing.T) {
	callCount := 0
	execCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
			},
		},
	}

	err := RegistryUpdate(cfg, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 0 {
		t.Errorf("expected 0 exec calls, got %d", callCount)
	}
}

func TestRegistryUpdateNoRegistry(t *testing.T) {
	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
		},
	}

	err := RegistryUpdate(cfg, false)
	if err == nil {
		t.Fatal("expected error for missing registry config")
	}
	if err.Error() != "no registry configured in deploy.yml" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegistryAddExecFailure(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name:   "community",
				Source: "https://example.com/registry",
			},
		},
	}

	err := RegistryAdd(cfg, false)
	if err == nil {
		t.Fatal("expected error from exec failure")
	}
}

func TestRegistryDeleteExecFailure(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
			},
		},
	}

	err := RegistryDelete(cfg, false)
	if err == nil {
		t.Fatal("expected error from exec failure")
	}
}

func TestRegistryUpdateExecFailure(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{
			Name: "my-pack",
			Registry: &config.RegistryConfig{
				Name: "community",
			},
		},
	}

	err := RegistryUpdate(cfg, false)
	if err == nil {
		t.Fatal("expected error from exec failure")
	}
}

func TestJobNameVarWins(t *testing.T) {
	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{
			Name: "from-deploy",
			Vars: map[string]string{"job_name": "from-var"},
		},
	}

	name, err := jobName(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "from-var" {
		t.Errorf("expected from-var, got %s", name)
	}
}

func TestJobNameFallbackToDeployName(t *testing.T) {
	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{
			Name: "from-deploy",
			Vars: map[string]string{"other": "value"},
		},
	}

	name, err := jobName(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "from-deploy" {
		t.Errorf("expected from-deploy, got %s", name)
	}
}

func TestJobNameMissing(t *testing.T) {
	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{
			Vars: map[string]string{"job_name": ""},
		},
	}

	_, err := jobName(cfg)
	if err == nil {
		t.Fatal("expected error when no job name can be determined")
	}
	if !strings.Contains(err.Error(), "cannot determine job name") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStatus(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{Name: "my-job"},
		Pack:   config.PackConfig{Name: "my-pack"},
	}

	if err := Status(cfg, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad", "job", "status", "my-job"}
	assertSliceEqual(t, expected, calls[0])
}

func TestStatusWithJobNameVar(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{
			Name: "my-job",
			Vars: map[string]string{"job_name": "var-job"},
		},
	}

	if err := Status(cfg, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad", "job", "status", "var-job"}
	assertSliceEqual(t, expected, calls[0])
}

func TestStatusWithExtraArgs(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{Name: "my-job"},
	}

	if err := Status(cfg, false, "-json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"nomad", "job", "status", "my-job", "-json"}
	assertSliceEqual(t, expected, calls[0])
}

func TestStatusDryRun(t *testing.T) {
	callCount := 0
	execCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{Name: "my-job"},
	}

	if err := Status(cfg, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 0 {
		t.Errorf("expected 0 exec calls, got %d", callCount)
	}
}

func TestStatusMissingJobName(t *testing.T) {
	callCount := 0
	execCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		return exec.Command("echo", "fake")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Pack: config.PackConfig{Name: "my-pack"},
	}

	err := Status(cfg, false)
	if err == nil {
		t.Fatal("expected error for missing job name")
	}
	if callCount != 0 {
		t.Errorf("expected 0 exec calls, got %d", callCount)
	}
}

func TestStatusExecFailure(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { execCommand = exec.Command }()

	cfg := &config.DeployConfig{
		Deploy: config.DeploySection{Name: "my-job"},
	}

	if err := Status(cfg, false); err == nil {
		t.Fatal("expected error from exec failure")
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
