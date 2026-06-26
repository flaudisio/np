// Package config reads and validates deploy.yml configuration files.
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// RegistryConfig configures a nomad-pack registry.
// When set, the registry is added automatically before running the pack.
type RegistryConfig struct {
	Name   string `mapstructure:"name"`
	Source string `mapstructure:"source"`
	Ref    string `mapstructure:"ref"`
}

// PackConfig identifies the nomad-pack pack and its optional registry.
type PackConfig struct {
	Name     string          `mapstructure:"name"`
	Registry *RegistryConfig `mapstructure:"registry"`
}

// DeploySection contains deployment parameters passed to nomad-pack.
type DeploySection struct {
	Name     string            `mapstructure:"name"`
	Vars     map[string]string `mapstructure:"vars"`
	VarFiles []string          `mapstructure:"var_files"`
}

// PlanConfig controls plan output behavior.
type PlanConfig struct {
	Verbose bool `mapstructure:"verbose"`
}

// DeployConfig is the root configuration structure parsed from deploy.yml.
type DeployConfig struct {
	Deploy DeploySection `mapstructure:"deploy"`
	Pack   PackConfig    `mapstructure:"pack"`
	Plan   PlanConfig    `mapstructure:"plan"`
}

// Load reads the YAML config at path, unmarshals it into a [DeployConfig],
// validates required fields, and applies defaults. It returns an error if
// the file cannot be read, the YAML is invalid, or required fields are
// missing.
//
// Defaults applied:
//   - deploy.var_files defaults to ["variables.hcl"] when not set.
func Load(path string) (*DeployConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg DeployConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Pack.Name == "" {
		return nil, errors.New("pack.name is required")
	}

	if cfg.Pack.Registry != nil {
		if cfg.Pack.Registry.Name == "" {
			return nil, errors.New("pack.registry.name is required when registry is set")
		}
		if cfg.Pack.Registry.Source == "" {
			return nil, errors.New("pack.registry.source is required when registry is set")
		}
	}

	if !v.IsSet("deploy.var_files") {
		cfg.Deploy.VarFiles = []string{"variables.hcl"}
	}

	return &cfg, nil
}
