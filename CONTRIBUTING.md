# Contributing

Thanks for contributing to `globalAI`.

## Local checklist

Before opening a pull request, run:

```bash
go test ./...
go build ./cmd/globalai
bash scripts/functional_smoke.sh
```

## Contribution guidelines

- Keep features small and explain the user-facing motivation in your PR.
- Add or update tests for every behavior change.
- Update `README.md` and relevant files in `docs/` when behavior or architecture changes.
- Prefer targeted changes over broad refactors.

## Development notes

- The project intentionally uses the standard library for the CLI and embedded viewer.
- The discovery layer is allowlist-based by design.
- The viewer is local-only and should remain safe to run on a personal workstation.
