# AGENTS.md

Go CLI that turns a YAML config (`deploy.yml`) into `nomad-pack` / `nomad` commands.

## Stack

- Go 1.26, managed via [mise](https://mise.jdx.dev/)
- CLI: [Cobra](https://github.com/spf13/cobra); Config: [Viper](https://github.com/spf13/viper)
- Lint: golangci-lint + pre-commit (yamllint `--strict`, shellcheck)

## Setup / common commands

```bash
mise install        # go, golangci-lint, nomad, nomad-pack, git-cliff, shellcheck, ...
mise run test       # go test ./...   (alias: t)
mise run lint       # pre-commit --all-files + golangci-lint-full   (alias: l)
mise run build      # go build -o np ./cmd/np, version injected via ldflags   (alias: b)
```

Single package / single test:

```bash
go test ./internal/nomadpack -run TestBuildCommandWithVars
```

## Run

```bash
go run ./cmd/np plan            # reads ./deploy.yml
go run ./cmd/np deploy          # alias: run
go run ./cmd/np plan --dry-run  # print commands, execute nothing
```

## Architecture

- `cmd/np/main.go` — Cobra wiring only; commands delegate to `run()` / `runStatus()`.
- `internal/config` — loads/validates `deploy.yml` via Viper; `pack.name` required.
- `internal/nomadpack` — builds arg slices and shells out to `nomad-pack` / `nomad`.
- `internal/log` — colored stderr logging (always ANSI, no TTY check).

Every subcommand maps to a nomad-pack action; `status` runs `nomad job status`.
Extra args after `--` are forwarded verbatim: `np render --dry-run -- --no-format`.

## Testing

- External commands are faked by swapping a package var — no real Nomad needed:
  `old := nomadpack.SetExecCommand(fn); defer nomadpack.SetExecCommand(old)`.
- CLI tests re-exec the test binary with `TEST_CLI=1` (see `cmd/np/main_test.go`);
  copy that pattern for new end-to-end command tests.

## Config defaults / behavior

- Config path defaults to `deploy.yml` (`-c/--config` overrides it).
- `deploy.var_files` defaults to `["variables.hcl"]` when unset; `var_files: []` disables it.
- `deploy.vars.job_name` wins over `deploy.name` when resolving the job for `status`.
- `nomad-pack plan` exit code 1 and `run` exit code 2 are treated as success.
- `-C/--cd` chdirs before running, in `PersistentPreRunE`.
- Logs go to stderr; stdout is left for command output.

## Commits

Use [Conventional Commits](https://www.conventionalcommits.org); `git-cliff`
parses them to build the changelog and choose the version bump (`feat` → minor,
`fix` → patch, breaking → major). Grouping is configured in `cliff.toml`.

## Releasing

```bash
mise run bump-version   # git-cliff bump + bump-my-version, commits and tags vX.Y.Z
mise run git-push       # push current branch and tags
```

Tags `v*` trigger GoReleaser (linux/darwin × amd64/arm64).
