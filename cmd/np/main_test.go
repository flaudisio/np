package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIHelp(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		os.Args = []string{"np", "--help"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIHelp")
	cmd.Env = append(os.Environ(), "TEST_CLI=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("process exited: %v\n%s", err, out)
	}
	if !contains(string(out), "deploy") {
		t.Errorf("expected deploy in help, got: %s", out)
	}
	if !contains(string(out), "plan") {
		t.Errorf("expected plan in help, got: %s", out)
	}
}

func TestCLIDeployMissingConfig(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		os.Args = []string{"np", "deploy", "--config", "/nonexistent/path.yaml"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIDeployMissingConfig")
	cmd.Env = append(os.Environ(), "TEST_CLI=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config")
	}
}

func TestCLIDeployWithConfig_DryRun(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		_ = os.Chdir(dir)
		os.Args = []string{"np", "deploy", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
`
	_ = os.WriteFile(filepath.Join(dir, "deploy.yaml"), []byte(yaml), 0o644)

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIDeployWithConfig_DryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !contains(string(out), "run complete") {
		t.Errorf("expected run complete message, got: %s", out)
	}
	if !contains(string(out), "nomad-pack run test-pack") {
		t.Errorf("expected nomad-pack command in output, got: %s", out)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
