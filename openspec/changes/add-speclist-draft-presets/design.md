## Context

The current draft path creates the same `Why`, `Context`, and `Proposed Requirements` sections for every workflow. That is workable, but it makes proposal and design review feel generic and pushes reviewers to reshape the draft immediately after generation.

## Goals / Non-Goals

**Goals:**
- Support a small typed preset set for draft generation.
- Keep presets grounded in the same retrieval bundle rather than branching into separate generation systems.
- Preserve preset choice on the draft payload and export metadata.

**Non-Goals:**
- Add user-defined templates or persistent preset management.
- Introduce LLM prompt variants or model orchestration for each preset.
- Overhaul export rendering beyond preserving the preset and using existing section fallbacks.

## Decisions

Add a typed `DraftPreset` value to the backend domain and draft payload.
Rationale: preset choice is part of the draft contract and should remain visible through review and export.

Use four presets: `general`, `proposal`, `design`, and `requirements`.
Rationale: these cover the next practical review modes without turning the workbench into a template marketplace.

Build preset sections from the same retrieval results with deterministic summarizers.
Rationale: this keeps the implementation simple, testable, and compatible with the current non-LLM draft generation flow.

Keep `general` as the normalization fallback.
Rationale: existing clients and exports should continue to work if preset input is missing or unknown.

## Risks / Trade-offs

[Preset structures are still approximate] -> Reviewers can already edit drafts before export, so presets only need to be strong starting points.

[Some exports expect legacy section names] -> Preserve compatibility by using section fallback lookups in export rendering.
