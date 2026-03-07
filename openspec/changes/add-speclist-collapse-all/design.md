## Context

The workbench already tracks collapse state per section. That makes bulk collapse a small extension of the same local state model rather than a new system. The main design question is where to place the controls so they are visible without crowding each section card.

## Goals / Non-Goals

**Goals:**
- Let reviewers collapse every current draft section with one action.
- Let reviewers expand every current draft section with one action.
- Keep the behavior local to the active draft review session.

**Non-Goals:**
- Persist bulk collapse state across page reloads.
- Add saved review layouts.
- Add partial-range collapse or section grouping.

## Decisions

Place bulk controls near the existing draft-level review actions.
Rationale: collapse-all and expand-all affect the whole draft, so they belong with reset and add-section controls instead of inside each section card.

Implement `Collapse all` by marking each current section index as collapsed.
Rationale: this keeps behavior explicit and avoids depending on implicit defaults.

Implement `Expand all` by clearing the collapse-state map.
Rationale: expanded is already the default state, so clearing the local map is the simplest representation.

Keep reset behavior unchanged so resetting the draft still returns to a fully expanded baseline.
Rationale: review reset should continue to mean a clean restart from the generated draft.

## Risks / Trade-offs

[Bulk collapse can hide too much context at once] -> `Expand all` remains adjacent, and each collapsed card still shows heading and section status summary.

[Bulk controls may become redundant for short drafts] -> The actions are lightweight and only appear when a draft exists.
