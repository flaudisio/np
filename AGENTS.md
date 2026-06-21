# AGENTS.md

Go CLI that reads a YAML config (`deploy.yaml`) and produces `nomad-pack` commands.

## Stack

- Go 1.26, managed via [mise](https://mise.jdx.dev/)
- CLI framework: [Cobra](https://github.com/spf13/cobra)
- Config: [Viper](https://github.com/spf13/viper)
- Lint: golangci-lint
- Pre-commit hooks: pre-commit (yamllint, shellcheck, golangci-lint)

## Setup

```bash
mise install  # Project dependencies
```

## Run

```bash
go run ./cmd/np deploy          # deploys from deploy.yaml
go run ./cmd/np plan            # plans from deploy.yaml
go run ./cmd/np plan --dry-run  # prints commands without executing
```

## Lint

```bash
mise run lint
```

## Build

```bash
mise run build
```

## Test

```bash
mise run test
```

## Conventions

- Go standard formatting (gofmt / goimports)
- `.editorconfig` is authoritative for indentation (tabs for Go, spaces for YAML/Markdown)
- Log messages go to stderr (via `internal/log`), not stdout
