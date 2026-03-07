---
area: docs
component: tooling-context
kind: proposal
status: active
tags:
  - openspec
  - docker
  - security
  - lsp
---

# Add Tooling Context

## Why

FastSpec is currently missing durable repo-local context for the tooling stack it keeps referencing: OpenSpec, Docker Compose, Traefik, CrowdSec, Trivy, Rust linting, and LSP-oriented workflows. That leaves future agents to rediscover the same ecosystem facts from scratch.

The repository also lacks an explicit note on what ideas are worth borrowing from `anomalyco/opencode` and which are just inspiration rather than requirements.

## What Changes

This change adds:

- richer project context in `openspec/config.yaml`
- durable docs for the container, security, and Rust developer-tooling stack
- a focused inspiration note for OpenCode design ideas
- a stable OpenSpec spec for maintaining tooling context docs

## Non-Goals

- implementing Docker, Traefik, CrowdSec, or Trivy in this repository
- building MCP servers or an agent runtime in this slice
- copying OpenCode behavior directly into FastSpec
