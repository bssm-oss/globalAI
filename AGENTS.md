# AGENTS.md

This file is for future human and AI contributors working on `globalAI`.

## Project purpose

`globalAI` is a local-first Go CLI that visualizes AI prompt and instruction files that are otherwise scattered across project and user-level configuration locations. The current product surface is intentionally narrow: one binary, one main workflow, and a deterministic discovery model.

## Core product rules

1. Keep the tool local-first. Do not add network-dependent product behavior for the viewer.
2. Prefer deterministic allowlists over broad filesystem crawling.
3. Avoid leaking unrelated secrets by scanning arbitrary directories.
4. Keep the runtime dependency footprint extremely small.
5. Preserve the safety property that prompt contents are rendered as plain text, not trusted HTML.

## Architecture boundaries

- `cmd/globalai/` should stay thin and only wire process-level concerns.
- `internal/cli/` owns argument parsing and command orchestration.
- `internal/source/` owns discovery policy and source metadata generation.
- `internal/viewer/` owns HTTP serving and embedded asset delivery.
- `internal/browser/` owns best-effort browser opening.

If a new feature does not fit one of those boundaries cleanly, discuss the boundary before adding code.

## Design constraints

- Prefer the Go standard library unless a new dependency materially reduces complexity.
- Do not add a JavaScript build pipeline unless the UI scope changes enough to justify it.
- Do not replace the allowlist with generic recursive scanning.
- Keep the server bound to loopback unless the product requirement explicitly changes.

## Testing expectations

Every change should keep three layers healthy:

1. unit tests for the changed package
2. full `go test ./...`
3. the functional smoke test script under `scripts/`

If you change the viewer contract, update both the Go tests and the smoke test.

## Documentation expectations

- Update `README.md` when user-facing behavior changes.
- Add a change record under `docs/changes/` for meaningful project changes.
- Update `docs/architecture.md` when package responsibilities or discovery rules change.

## Git and release expectations

- Prefer small, reviewable commits grouped by concern.
- Keep test changes close to the implementation they validate.
- Do not ship changes with failing tests, failing builds, or broken smoke tests.
