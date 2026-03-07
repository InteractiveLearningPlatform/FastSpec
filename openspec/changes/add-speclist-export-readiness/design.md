## Context

The current workbench already has enough client-side information to estimate export readiness: the edited draft, the original draft snapshot, and the section review flags. The cheapest useful improvement is to compute readiness locally and show it before export instead of adding a backend validation endpoint.

## Goals / Non-Goals

**Goals:**
- Compute a simple readiness state from the current draft and review flags.
- Distinguish blocking issues from warnings.
- Show readiness status next to export controls and in one compact summary.

**Non-Goals:**
- Block export at the transport layer.
- Add backend readiness APIs or persistent review workflows.
- Replace existing backend export validation.

## Decisions

Compute readiness entirely in the React workbench.
Rationale: the draft and review state already live client-side, and backend export validation remains the final guardrail.

Treat `blocked` review flags and missing required draft fields as blockers.
Rationale: these conditions mean the reviewer has explicitly or structurally indicated the draft is not ready.

Treat `needs-work`, missing citations, and unchanged generated drafts as warnings.
Rationale: these conditions are useful review signals but should not be as severe as blockers.

Show a high-level readiness badge plus flat lists of blockers and warnings.
Rationale: reviewers need a fast go/no-go summary, not a complex report.

## Risks / Trade-offs

[Client-side readiness is advisory] -> Keep backend export validation as the final correctness layer.

[Some teams may want stricter gating] -> Keep this slice informational and defer hard gating to a later workflow change if needed.
