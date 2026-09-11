# np

CLI for deploying Nomad Pack applications from a `deploy.yml` config file.

## Install

### From GitHub releases

Download the latest `np` binary for your platform from the
[releases page](https://github.com/flaudisio/np/releases) and extract the archive.

### From source

```bash
go install github.com/flaudisio/np/cmd/np@latest
```

Or build locally with [mise](https://mise.jdx.dev/):

```bash
mise install
mise run build
```

## Quick Start

```bash
mise install
go run ./cmd/np plan
go run ./cmd/np deploy
```

## Configuration

`np` reads a single `deploy.yml` file. Point it at a different path with `--config`.

```yaml
pack:
  name: my-pack              # required
  registry:                  # optional
    name: my-registry
    source: github.com/org/pack-registry
    ref: main
deploy:
  name: my-deployment        # optional
  vars:                      # optional
    region: us-east-1
  var_files: [variables.hcl] # defaults to this when unset
plan:
  verbose: false             # optional
```

## Commands

| Command                    | Description                       |
| -------------------------- | --------------------------------- |
| `np deploy` (alias `run`)  | Deploy the pack                   |
| `np plan`                  | Plan the deployment               |
| `np destroy`               | Destroy the deployment            |
| `np stop`                  | Stop the deployment               |
| `np render`                | Render the pack                   |
| `np registry add`          | Add the configured registry       |
| `np registry delete`       | Delete the configured registry    |
| `np registry update`       | Update the configured registry    |

### Flags

| Flag               | Description                                  |
| ------------------ | -------------------------------------------- |
| `-c, --config`     | Path to `deploy.yml` (default: `deploy.yml`) |
| `-n, --dry-run`    | Print commands without executing             |
| `-C, --cd`         | Change directory before running commands     |

## License

[MIT](LICENSE)
