# Registry Subcommands Design

**Date:** 2026-07-02

## Purpose

Add `np registry add`, `np registry delete`, and `np registry update`
subcommands as config-driven wrappers around the underlying `nomad-pack
registry` commands.

## CLI structure

```
np registry add     → nomad-pack registry add <name> <source> [--ref <ref>]
np registry delete  → nomad-pack registry delete <name> [--ref <ref>]
np registry update  → nomad-pack registry update <name> [--ref <ref>]
```

`np registry` is a group command with no action of its own. All three
subcommands inherit the root's persistent flags (`--config`, `--dry-run`,
`--cd`). No positional args or local flags — all values come from
`deploy.yml`.

## Config

Reads from `pack.registry` in `deploy.yml`:

```yaml
pack:
  registry:
    name: community
    source: github.com/hashicorp/nomad-pack-community-registry
    ref: v0.1.0   # optional
```

Validation:
- `add`: `name` and `source` required (enforced by existing `config.Load`)
- `delete`: `name` required (validated at call time)
- `update`: `name` required (validated at call time)

## New functions in `internal/nomadpack`

Three exported functions, each following the existing `ensureRegistry`
pattern (build args, log command, dry-run check, execute via `execCommand`):

```go
func RegistryAdd(cfg *config.DeployConfig, dryRun bool) error
func RegistryDelete(cfg *config.DeployConfig, dryRun bool) error
func RegistryUpdate(cfg *config.DeployConfig, dryRun bool) error
```

Each validates that `cfg.Pack.Registry` is non-nil and has the required
fields, constructs the arg slice, logs, and either returns (dry-run) or
executes.

## CLI wiring (`cmd/np/main.go`)

Four factory functions: `registryCmd()` (group), `registryAddCmd()`,
`registryDeleteCmd()`, `registryUpdateCmd()`. Each subcommand's `RunE`
loads config then calls the corresponding `nomadpack` function.

No reuse of the shared `run()` function — registry commands don't need
cwd logging, action completion messages, or the `nomadpack.Run` codepath.

## Error handling

- Missing `pack.registry` in config → `"no registry configured in
  deploy.yml"`
- Missing required fields → errors from config or runtime validation
- Exec errors propagate via Cobra's `RunE`, printed to stderr

## What doesn't change

- `config.Load` — registry validation already present
- `ensureRegistry` — unchanged, still used as pre-flight for
  deploy/plan/etc
- Existing subcommands and `run()` — untouched
- No `--target` support (out of scope)

## Testing

- **Unit (`nomadpack_test.go`)**: Use `execCommand` mock pattern. Verify
  args, dry-run skip, missing-config errors, exec failure propagation
- **Integration (`main_test.go`)**: Follow existing subprocess testing
  pattern (`TEST_CLI` env var, temp dirs). Test dry-run output and error
  messages
