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

func TestCLIDestroyDryRun(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "destroy", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIDestroyDryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "destroy complete") {
		t.Errorf("expected destroy complete message, got: %s", out)
	}
	if !strings.Contains(string(out), "nomad-pack destroy test-pack") {
		t.Errorf("expected nomad-pack destroy command in output, got: %s", out)
	}
}

func TestCLIStopDryRun(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "stop", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIStopDryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "stop complete") {
		t.Errorf("expected stop complete message, got: %s", out)
	}
	if !strings.Contains(string(out), "nomad-pack stop test-pack") {
		t.Errorf("expected nomad-pack stop command in output, got: %s", out)
	}
}

func TestCLIRenderDryRun(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "render", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIRenderDryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "render complete") {
		t.Errorf("expected render complete message, got: %s", out)
	}
	if !strings.Contains(string(out), "nomad-pack render test-pack") {
		t.Errorf("expected nomad-pack render command in output, got: %s", out)
	}
}

func TestCLIRenderDryRunWithExtraArgs(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "render", "--dry-run", "--", "--no-format", "--arg", "value"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIRenderDryRunWithExtraArgs")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "nomad-pack render test-pack") {
		t.Errorf("expected nomad-pack render command, got: %s", out)
	}
	if !strings.Contains(string(out), "--no-format") {
		t.Errorf("expected --no-format in output, got: %s", out)
	}
	if !strings.Contains(string(out), "--arg value") {
		t.Errorf("expected --arg value in output, got: %s", out)
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
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
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

func TestCLIRegistryAddDryRun(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "registry", "add", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
  registry:
    name: community
    source: https://example.com/registry
    ref: v0.1.0
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIRegistryAddDryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "nomad-pack registry add community https://example.com/registry --ref v0.1.0") {
		t.Errorf("expected full registry add command in output, got: %s", out)
	}
}

func TestCLIRegistryDeleteDryRun(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "registry", "delete", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
  registry:
    name: community
    source: https://example.com/registry
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIRegistryDeleteDryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "nomad-pack registry delete community") {
		t.Errorf("expected registry delete command in output, got: %s", out)
	}
}

func TestCLIRegistryUpdateDryRun(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "registry", "update", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
  registry:
    name: community
    source: https://example.com/registry
    ref: v0.2.0
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIRegistryUpdateDryRun")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "nomad-pack registry update community --ref v0.2.0") {
		t.Errorf("expected full registry update command in output, got: %s", out)
	}
}

func TestCLIRegistryAddMissingConfig(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		os.Args = []string{"np", "registry", "add", "--config", "/nonexistent/path.yaml"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIRegistryAddMissingConfig")
	cmd.Env = append(os.Environ(), "TEST_CLI=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config")
	}
}

func TestCLIRegistryAddNoRegistry(t *testing.T) {
	if os.Getenv("TEST_CLI") == "1" {
		dir := os.Getenv("TEST_DIR")
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"np", "registry", "add", "--dry-run"}
		main()
		return
	}

	dir := t.TempDir()
	yaml := `
pack:
  name: test-pack
`
	if err := os.WriteFile(filepath.Join(dir, "deploy.yml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCLIRegistryAddNoRegistry")
	cmd.Env = append(os.Environ(), "TEST_CLI=1", "TEST_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for missing registry")
	}
	if !strings.Contains(string(out), "no registry configured") {
		t.Errorf("expected 'no registry configured' error, got: %s", out)
	}
}
