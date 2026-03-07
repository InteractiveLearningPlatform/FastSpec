---
area: docs
component: inspiration
kind: reference
status: active
language: multiple
framework: opencode
tags:
  - opencode
  - agents
  - lsp
---

# OpenCode Ideas

This document captures ideas worth considering from `anomalyco/opencode`. It is intentionally about inspiration, not current FastSpec behavior.

## Useful Design Ideas

### Provider-Agnostic Execution

OpenCode emphasizes that it is not coupled to a single model provider and can work with Claude, OpenAI, Google, or local models. That is a strong fit for FastSpec as well.

Implication for FastSpec:

- avoid baking provider-specific assumptions into the durable spec model
- keep model and tool routing as pluggable runtime concerns

### Explicit Plan And Build Modes

OpenCode exposes two built-in agents:

- `build` for full-access development work
- `plan` for read-only analysis and exploration

That separation is valuable. FastSpec should likely distinguish planning artifacts from execution artifacts and may later benefit from agent modes with different permissions.

### Subagents For Focused Work

OpenCode includes a general subagent for complex searches and multistep tasks.

FastSpec implication:

- durable specs should be composable enough that a coordinator can hand a narrow slice to a helper agent without copying the entire project context

### Out-Of-The-Box LSP Support

OpenCode highlights out-of-the-box LSP support as a core differentiator.

FastSpec implication:

- semantic code intelligence should be considered part of the expected agent environment
- future runtime design should be comfortable with LSP-backed indexing and navigation

### Strong Terminal UX

OpenCode has a strong terminal-first bias and treats the TUI as a serious product surface.

FastSpec implication:

- if a CLI or TUI is built, it should not feel like a thin debug shell
- agent-oriented systems still need a deliberate operator experience

### Client/Server Architecture

OpenCode notes that its TUI is only one client and that a client/server architecture enables remote control from other surfaces.

FastSpec implication:

- keep protocol and model layers separate from UI surfaces
- avoid designs that trap FastSpec inside a single editor or terminal frontend

## What Not To Infer

These notes do not mean FastSpec currently has:

- multiple agent modes
- a subagent runtime
- provider routing
- client/server architecture

They only identify design directions worth preserving as the project evolves.
