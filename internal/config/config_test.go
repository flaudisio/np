package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "deploy.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMinimal(t *testing.T) {
	yaml := `
pack:
  name: my-pack
`
	path := writeTempYAML(t, yaml)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pack.Name != "my-pack" {
		t.Errorf("expected my-pack, got %s", cfg.Pack.Name)
	}
}

func TestLoadFull(t *testing.T) {
	yaml := `
deploy:
  name: my-deploy
  vars:
    key1: value1
    key2: "42"
  var_files:
    - vars/common.yaml

pack:
  name: my-pack
  registry:
    name: my-registry
    source: https://github.com/org/repo
    ref: v1.2.3

plan:
  verbose: true
`
	path := writeTempYAML(t, yaml)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Deploy.Name != "my-deploy" {
		t.Errorf("expected my-deploy, got %s", cfg.Deploy.Name)
	}
	if cfg.Deploy.Vars["key1"] != "value1" {
		t.Errorf("expected value1, got %s", cfg.Deploy.Vars["key1"])
	}
	if cfg.Deploy.Vars["key2"] != "42" {
		t.Errorf("expected 42, got %s", cfg.Deploy.Vars["key2"])
	}
	if len(cfg.Deploy.VarFiles) != 1 || cfg.Deploy.VarFiles[0] != "vars/common.yaml" {
		t.Errorf("expected vars/common.yaml, got %v", cfg.Deploy.VarFiles)
	}
	if cfg.Pack.Name != "my-pack" {
		t.Errorf("expected my-pack, got %s", cfg.Pack.Name)
	}
	if cfg.Pack.Registry.Name != "my-registry" {
		t.Errorf("expected my-registry, got %s", cfg.Pack.Registry.Name)
	}
	if cfg.Pack.Registry.Source != "https://github.com/org/repo" {
		t.Errorf("expected source, got %s", cfg.Pack.Registry.Source)
	}
	if cfg.Pack.Registry.Ref != "v1.2.3" {
		t.Errorf("expected v1.2.3, got %s", cfg.Pack.Registry.Ref)
	}
	if !cfg.Plan.Verbose {
		t.Error("expected verbose to be true")
	}
}

func TestDefaults(t *testing.T) {
	yaml := `
pack:
  name: my-pack
`
	path := writeTempYAML(t, yaml)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Deploy.Name != "" {
		t.Errorf("expected empty deploy name, got %s", cfg.Deploy.Name)
	}
	if len(cfg.Deploy.Vars) != 0 {
		t.Errorf("expected empty vars, got %v", cfg.Deploy.Vars)
	}
	if len(cfg.Deploy.VarFiles) != 0 {
		t.Errorf("expected empty var_files, got %v", cfg.Deploy.VarFiles)
	}
	if cfg.Plan.Verbose {
		t.Error("expected verbose to be false by default")
	}
	if cfg.Pack.Registry != nil {
		t.Error("expected registry to be nil by default")
	}
}

func TestMissingPackName(t *testing.T) {
	yaml := `
deploy:
  name: foo
`
	path := writeTempYAML(t, yaml)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing pack.name")
	}
	if !strings.Contains(err.Error(), "pack.name is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegistryMissingName(t *testing.T) {
	yaml := `
pack:
  name: my-pack
  registry:
    source: https://example.com
`
	path := writeTempYAML(t, yaml)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing registry name")
	}
	if !strings.Contains(err.Error(), "pack.registry.name is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegistryMissingSource(t *testing.T) {
	yaml := `
pack:
  name: my-pack
  registry:
    name: my-registry
`
	path := writeTempYAML(t, yaml)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing registry source")
	}
	if !strings.Contains(err.Error(), "pack.registry.source is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
