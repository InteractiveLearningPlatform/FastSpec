## Context

The current FastSpec export path writes a generic `FastSpecDraft` YAML document regardless of what the draft contains. That is serviceable as a transport format, but it does not give users a strong starting point for durable spec refinement the way typed OpenSpec templates now do for markdown exports.

## Goals / Non-Goals

**Goals:**
- Render a more structured FastSpec YAML shape for exported drafts.
- Preserve citations and draft summary metadata.
- Improve the YAML output without introducing brittle inference.

**Non-Goals:**
- Perfectly infer all possible FastSpec document kinds.
- Replace later manual refinement of exported YAML.
- Add UI changes if the export contract remains compatible.

## Decisions

Render a typed `SpecDocumentDraft`-style YAML structure instead of the existing generic `FastSpecDraft`.
Rationale: the export should look more like a purposeful durable draft document, not a raw transport wrapper.

Map common draft content into stable sections such as rationale, context, and requirements.
Rationale: these are reusable across many durable spec cases without forcing over-specific type inference.

Keep citation sidecars unchanged.
Rationale: provenance handling is already correct and should not be coupled to the primary YAML shape.

## Risks / Trade-offs

[The YAML remains somewhat generic] -> Prefer a stable typed draft document over risky over-inference.

[Users may expect full project/module inference] -> Keep this slice scoped to stronger draft structure and defer precise type inference to later work.
