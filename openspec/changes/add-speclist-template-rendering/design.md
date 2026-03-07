## Context

Repo-aware change targeting solved the path problem for Speclist exports, but not the content-shape problem. A draft exported to `proposal.md` should not look identical to a draft exported to `tasks.md` or `specs/<capability>/spec.md`. This slice introduces template-aware rendering while still using the reviewed draft object as the source of truth.

## Goals / Non-Goals

**Goals:**
- Render different OpenSpec artifact classes with different output templates.
- Map reviewed draft sections into a more appropriate artifact shape.
- Keep citation visibility in rendered content.

**Non-Goals:**
- Fully infer perfect artifact semantics from arbitrary drafts.
- Replace later human editing of exported artifacts.
- Add schema-specific logic for every future OpenSpec schema variant.

## Decisions

Render proposal, design, tasks, and spec targets with separate template functions.
Rationale: the artifact shapes are distinct enough that one generic markdown renderer is no longer sufficient.

Keep filesystem markdown export on the generic renderer for now.
Rationale: the typed shape is specifically valuable for OpenSpec artifact targets; generic markdown export can stay simple in this slice.

Generate conservative task checklists and requirement scaffolds from draft sections.
Rationale: this produces useful starting structure without pretending the draft can infer every checkbox or scenario perfectly.

## Risks / Trade-offs

[Rendered templates may guess the wrong section mapping] -> Use deterministic heuristics and keep exports clearly editable.

[Overfitting to current OpenSpec schema] -> Limit this slice to the currently supported proposal/design/tasks/spec targets.
