# globalAI

`globalAI` is a Go CLI for inspecting AI prompt and instruction sources that are usually hard to review in one place. The first shipped workflow is `globalai web`, which starts a local viewer for curated project and global instruction files such as `AGENTS.md`, `CLAUDE.md`, `.claude/`, `.cursor/rules/`, and `.sisyphus/`.

The project is intentionally dependency-light. The CLI uses the Go standard library for command handling, discovery, the local HTTP server, and embedded static assets. That keeps the binary simple to audit, simple to ship, and easy for future contributors to understand.

## Why this exists

Many AI tools store important instructions in different places: repository files, global user config folders, and tool-specific rule directories. Those files become hard to audit as the number of tools grows. `globalAI` gives you one local viewer that makes those sources visible without introducing a frontend toolchain or a network dependency.

## Current feature set

- `globalai web` starts a loopback-only viewer.
- The viewer exposes a curated, allowlisted set of source locations instead of scanning the entire filesystem.
- Prompt contents are rendered as plain text in the browser to avoid HTML injection from raw source files.
- The embedded UI is shipped directly inside the Go binary.

## Discovery model

The current allowlist favors clarity and safety over broad scanning.

### Project-local sources

- `AGENTS.md`
- `CLAUDE.md`
- `GEMINI.md`
- `.github/copilot-instructions.md`
- `.claude/**`
- `.cursor/rules/**`
- `.sisyphus/**`

### Global sources

- `~/AGENTS.md`
- `~/.claude/**`
- `~/.cursor/rules/**`
- `~/.sisyphus/**`

Supported file extensions are currently `.md`, `.mdc`, `.txt`, `.json`, `.yaml`, and `.yml`. Files larger than 1 MiB are skipped.

## Getting started

### Requirements

- Go 1.25+

### Build

```bash
go build ./cmd/globalai
```

### Run the viewer

```bash
./globalai web
```

You can also target a specific root explicitly.

```bash
./globalai web --root /path/to/repository
./globalai web /path/to/repository
```

### Useful flags

- `--addr`: override the listening address, default `127.0.0.1:0`
- `--open`: force browser opening
- `--no-open`: skip browser opening, useful for CI and smoke tests

## Development workflow

### Format, test, and build

```bash
gofmt -w $(find . -name '*.go' -not -path './.git/*')
go test ./...
go build ./cmd/globalai
```

### Functional smoke test

```bash
bash scripts/functional_smoke.sh
```

That script builds the binary, starts `globalai web --no-open`, probes the local API, and verifies the returned JSON.

## CI

The repository ships with a GitHub Actions workflow that runs:

- formatting verification via `gofmt -s`
- `go vet`
- `go test -race ./...`
- `go build ./cmd/globalai`
- the functional smoke test script

## Project structure

```text
cmd/globalai/           CLI entrypoint
internal/cli/           command parsing and orchestration
internal/browser/       OS browser launcher
internal/source/        allowlisted prompt source discovery
internal/viewer/        HTTP server and embedded static assets
docs/                   architecture notes and change records
scripts/                repeatable verification helpers
```

## Documentation for contributors and AI agents

- `AGENTS.md` explains repository rules, architecture boundaries, and expectations for future automated edits.
- `docs/architecture.md` explains the initial system shape.
- `docs/changes/2026-04-06-initial-bootstrap.md` records the first bootstrap change in a durable format.

## Roadmap

- richer source family support for more AI tools
- search and filtering inside the viewer
- packaged releases for macOS, Linux, and Windows
- cleaner export or snapshot workflows for auditing prompt state over time
