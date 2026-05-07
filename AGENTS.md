# AGENTS.md

Always reply in Chinese.
This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## Project Overview

**prompthub** — a prompt management/sharing platform. Early stage (initial scaffolding).

- **Remote**: `git@github.com:mymikasa/prompthub.git`
- **License**: MIT (Copyright 2026 Mikasa)

## Structure

- `backend/` — Go backend (empty, not yet scaffolded)
- `frontend/` — Frontend application (empty, not yet scaffolded)

## Go Conventions

The `.gitignore` is Go-oriented. When the backend is scaffolded:
- Test binaries: `*.test`, coverage: `*.out`, `coverage.*`
- Workspace files: `go.work`, `go.work.sum` are gitignored
- Environment: `.env` is gitignored
