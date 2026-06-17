package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	if !strings.Contains(string(out), "deploy") {
		t.Errorf("expected deploy in help, got: %s", out)
	}
	if !strings.Contains(string(out), "plan") {
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
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "deploy", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIDeployWithConfig_DryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "run complete") {
		t.Errorf("expected run complete message, got: %s", out)
	}
	if !strings.Contains(string(out), "nomad-pack run test-pack") {
		t.Errorf("expected nomad-pack command in output, got: %s", out)
	}
}
