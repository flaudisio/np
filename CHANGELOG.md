# Changelog

## Unreleased

### Features

- Add status command

## v0.1.0 - 2026-09-11

### Features

- Add --version flag with ldflags injection
- Add 'reg' alias to registry command
- Add registry add/delete/update CLI subcommands
- Add registry add/delete/update functions to nomadpack
- Add --cd/-C flag for working directory override
- Add debug log level and gate 'complete' messages behind --debug flag
- Support arbitrary extra args via -- separator
- Add render subcommand
- Add destroy and stop subcommands
- Add 'run' as alias for 'deploy' subcommand
- Default var_files to ['variables.hcl'] when not defined
- Add Cobra CLI with deploy and plan commands
- Add Viper config loading and validation
- Add ANSI stderr logger
- Scaffold np project structure

### Bug fixes

- [release] Build from cmd/np main package
- Add field validation and exec failure tests for registry commands
- Treat nomad-pack run exit code 2 as non-error
- Treat nomad-pack plan exit code 1 as non-error

### Refactor

- Change success messag to info and remove unused debug logic
- Use `deploy.yml` as default config file

### Documentation

- Add config reference, commands, install, and releasing notes
- Cleanup old specs
- Add design doc for registry subcommands
- Add Godoc comments to all exported symbols across all packages
- Add AGENTS.md and README

### Styling

- [release] Use double quotes consistently

### Miscellaneous

- [release] Add git-cliff and bump-my-version setup
- [gitignore] Align with toptal scaffold template
- [deps] Configure dependabot for go and actions
- [.github] Update pre-commit workflow
- Add release pipeline with goreleaser
- [ci] Add pre-commit workflow
- [pre-commit] Update hooks
- [mise] Bump golangci-lint to 2.12
- [mise] Add Nomad tools
- [mise] Update tasks
- Add license file
