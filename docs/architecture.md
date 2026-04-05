# Architecture

## Overview

`globalAI` is a greenfield Go CLI built around one core workflow: discover AI prompt and instruction files, then show them in a local embedded web viewer.

The architecture is intentionally small.

## Runtime flow

1. `cmd/globalai/main.go` creates a signal-aware context.
2. `internal/cli` parses command-line arguments and resolves root and home directories.
3. `internal/source` collects allowlisted files from project and global locations.
4. `internal/viewer` starts a loopback HTTP server and exposes both the embedded UI and `/api/sources`.
5. `internal/browser` optionally opens the local viewer URL with the platform default browser.

## Why the project is structured this way

### Standard library first

The current product scope does not justify Cobra, Viper, a JS build chain, or a template/rendering dependency. Using stdlib-only plumbing keeps the codebase transparent and lowers maintenance cost for an OSS project that is just getting started.

### Discovery is deterministic

The tool does not scan the entire filesystem. It only checks a curated set of known prompt and rule locations. That choice is product-driven:

- it reduces accidental exposure of unrelated files
- it keeps viewer behavior predictable
- it makes documentation and tests much easier to keep accurate

The current supported matrix is intentionally narrow:

- project root files: `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`
- project config locations: `.github/copilot-instructions.md`, `.claude/**`, `.cursor/rules/**`, `.sisyphus/**`
- global locations: `~/AGENTS.md`, `~/.claude/**`, `~/.cursor/rules/**`, `~/.sisyphus/**`

Files outside those locations are not part of the current product contract.

### Embedded assets

The viewer is shipped via `go:embed`, which means contributors do not need a separate frontend toolchain to work on the project. Static assets live under `internal/viewer/static/` and are served directly by the binary.

The viewer is loopback-only and the browser UI renders both prompt bodies and filesystem-derived metadata through text-safe DOM updates instead of trusting those values as HTML.

## Current package map

- `internal/cli`: command entrypoints and orchestration
- `internal/source`: discovery rules and payload generation
- `internal/viewer`: HTTP serving and session lifecycle
- `internal/browser`: OS-specific browser opening

## Near-term extension points

- richer discovery families for more AI tools
- filtering and search in the embedded UI
- release automation for tagged builds

Any larger feature should preserve the current safety guarantees unless the product explicitly decides otherwise.
