package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type RegistryConfig struct {
	Name   string `mapstructure:"name"`
	Source string `mapstructure:"source"`
	Ref    string `mapstructure:"ref"`
}

type PackConfig struct {
	Name     string          `mapstructure:"name"`
	Registry *RegistryConfig `mapstructure:"registry"`
}

type DeploySection struct {
	Name     string            `mapstructure:"name"`
	Vars     map[string]string `mapstructure:"vars"`
	VarFiles []string          `mapstructure:"var_files"`
}

type PlanConfig struct {
	Verbose bool `mapstructure:"verbose"`
}

type DeployConfig struct {
	Deploy DeploySection `mapstructure:"deploy"`
	Pack   PackConfig    `mapstructure:"pack"`
	Plan   PlanConfig    `mapstructure:"plan"`
}

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
