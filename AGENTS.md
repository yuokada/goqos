# Repository Guidelines

## Project Structure & Module Organization
`main.go` hosts the CLI entry point, handles the `set`/`view`/`clear` subcommand flow, and injects the build metadata. The helper script `bin/qos.sh` calls `tc` and related system tools on behalf of the binary. Service boot artifacts live under `etc/rc.d/init.d` and `etc/sysconfig` for CentOS-style deployments. RPM packaging assets are grouped in `buildrpm.sh` and `rpm/goqos.spec`, while the `Dockerfile` recreates the CentOS 7 toolchain required for reproducible RPM builds.

## Build, Test, and Development Commands
- `make build`: Fetches dependencies if needed and produces the binary with tag and commit embedded via `-ldflags`.
- `make build_linux`: Cross-compiles the CLI for Linux/amd64 from macOS or other hosts.
- `make rpm`: Runs `buildrpm.sh` to assemble RPM outputs inside `rpm/` and `rpmbuild/`.
- `go run .` / `./goqos set ...`: Exercise the CLI locally; combine with `view` and `clear` to verify end-to-end flows.

## Coding Style & Naming Conventions
Always run `go fmt ./...` (or `gofmt`) before committing. Exported identifiers use PascalCase, internal identifiers use camelCase, and errors should be wrapped with `github.com/pkg/errors` to keep stack context. Shell scripts and spec files should follow POSIX syntax with consistent two-space indentation for YAML fragments.

## Testing Guidelines
There are no `_test.go` files yet, so accompany new logic with table-driven tests and keep `go test ./...` green. Abstract system-dependent behavior behind interfaces so you can mock `tc` calls in unit tests. At minimum, cover representative `set` and `clear` scenarios to prevent regressions in traffic control rules.

## Commit & Pull Request Guidelines
Recent history favors concise English imperative subjects such as "Create ..." or "Update ...". Keep the subject within 50 characters and append references like `Refs #123` when linking issues. Pull requests should describe the intent, document test evidence (`go test ./...`, `make build`, etc.), note RPM or Docker considerations, and include sample CLI output (e.g., `./goqos view`) when behavior changes.

## Release & Configuration Tips
For RPM distribution, rely on the `Dockerfile` build environment and collect artifacts from `rpmbuild/RPMS/`. When installing on CentOS, place `etc/rc.d/init.d/goqos` under `/etc/init.d` and adjust `/etc/sysconfig/goqos` to match the host. Because traffic shaping requires elevated privileges, avoid running these commands in CI; favor staged manual verification instead.
