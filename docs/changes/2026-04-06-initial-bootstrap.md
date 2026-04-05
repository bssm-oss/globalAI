# 2026-04-06 — Initial bootstrap

## Summary

This change bootstraps the `globalAI` repository from an empty state into a working Go CLI project with a local web viewer, tests, documentation, and CI automation.

## What shipped

- a `globalai` Go binary with the `web` command
- deterministic discovery for project and global AI instruction sources
- a local embedded viewer served over loopback with an `/api/sources` endpoint
- unit tests for CLI behavior, discovery, browser launching, and viewer serving
- a functional smoke test script that verifies the built binary end to end
- repository docs including `README.md`, `AGENTS.md`, `CONTRIBUTING.md`, and `docs/architecture.md`
- GitHub Actions CI for formatting, vetting, tests, build, and smoke verification

## Product decisions

- standard library first
- no frontend build toolchain
- allowlist-based discovery instead of broad recursive scanning
- text-only source rendering to avoid trusting raw prompt file contents as HTML

## Follow-up ideas

- add more AI tool source families once real user demand is clear
- introduce release automation when the command surface stabilizes
- add viewer-side search and filtering once the source catalog grows
