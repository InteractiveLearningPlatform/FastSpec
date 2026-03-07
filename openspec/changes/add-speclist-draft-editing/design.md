## Context

The current workbench shows draft sections and citations but treats them as fixed output. Since export already takes the full draft payload from the frontend, the cheapest useful improvement is to let reviewers edit that payload in the browser and then validate it server-side before durable export.

## Goals / Non-Goals

**Goals:**
- Make draft title, summary, sections, and citations editable in the workbench.
- Allow reviewers to add and remove sections.
- Normalize edited drafts before export and reject invalid structures.
- Keep the implementation local to the current product/workbench slice.

**Non-Goals:**
- Add persistent draft storage, autosave, or collaboration.
- Introduce a separate draft-editing backend API.
- Build rich-text editing or a full document editor.

## Decisions

Keep draft editing entirely client-side until export.
Rationale: the export request already carries the full draft object, so this adds useful editing with minimal architectural surface area.

Treat sections as plain text fields with editable citation lists.
Rationale: markdown and YAML export already consume plain strings, and a textarea-based editor is enough for this slice.

Normalize and validate edited drafts during export.
Rationale: reviewers should be able to clean up content freely, but the backend still needs to reject empty titles, empty sections, or blank headings before writing artifacts.

Add explicit section add/remove controls.
Rationale: reviewers need a lightweight way to reshape generated drafts, not just tweak existing text.

## Risks / Trade-offs

[Client-side editing can diverge from generated defaults] -> Export validation remains the backend guardrail.

[Textarea editing is not rich enough for all future workflows] -> Keep this slice intentionally simple and defer richer editing to later product changes.

[Manual citation edits could reduce provenance quality] -> Keep citations editable but visible as individual lines so reviewers can preserve or correct them intentionally.
