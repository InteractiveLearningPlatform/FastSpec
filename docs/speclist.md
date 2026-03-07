# Speclist

Speclist is the ingestion and retrieval product surface in this repo.

The current product direction extends Speclist into a broader marketplace and
platform. The operating model, storage topology, security gates, and delivery
baseline for that direction are defined in `docs/speclist-platform-ops.md`.

It exists to bridge the gap between:

- existing documentation that teams already have in DOCX or Confluence
- durable FastSpec and OpenSpec artifacts that agents can reason over efficiently

## Components

- `apps/speclist-api/`
  Go backend microservice using hexagonal architecture.
- `apps/speclist-web/`
  React frontend for operators and reviewers.

## Workflow

1. Import DOCX or Confluence content into the Speclist corpus.
2. Index existing FastSpec and OpenSpec repository artifacts.
3. Search for grounded context bundles instead of raw documents.
4. Generate reviewable draft specs with source citations.
5. Review and edit draft titles, summaries, sections, and citations in the workbench.
6. Inspect any citation to reopen the grounded source chunk behind it.
7. Inspect indexed sources directly to review source metadata and chunk inventories.
8. Choose a draft preset to start proposal, design, or requirements-oriented review from a better section structure.
9. Compare the current edited draft against the original generated draft before export.
10. Mark draft sections with lightweight review flags and notes before export.
11. Review a compact export-readiness summary before export.
12. Reset the current review back to the original generated draft when needed.
13. Reorder draft sections during review without recreating them.
14. Duplicate draft sections during review to branch local variations quickly.
15. Export accepted drafts into durable OpenSpec markdown or FastSpec YAML files.
16. Optionally target an active OpenSpec change artifact directly instead of a generic output path.
17. Narrow retrieval and drafting with source-kind, source-origin, and location filters when the corpus mixes imported docs and repository specs.

## Design Constraints

- Keep retrieval compact and interactive.
- Preserve provenance for every returned excerpt.
- Keep the backend domain independent from storage, HTTP, and source-specific adapters.
- Treat generated drafts as reviewable candidates, not final truth.
- Require explicit export destinations and avoid silent overwrite of existing files.
- Keep the workbench as one platform surface, not the full product boundary.
- Keep retrieval filters simple and source-oriented so they stay compatible with later platform-scale search work.
- Let reviewers refine generated drafts before export, while keeping export validation on the backend.
- Keep citation verification close to the draft review workflow.
- Keep source-level corpus review available from the indexed source list.
- Let draft generation start from a small set of intentional review presets.
- Keep edit review explicit by showing how the current draft differs from the generated draft.
- Keep section-level review concerns visible during the final review pass.
- Keep export decisions informed by a compact readiness summary.
- Let reviewers restart from the generated baseline without regenerating a draft.
- Let reviewers refine section order during review without rebuilding the draft.
- Let reviewers branch an existing section into a local variation without copy-paste.

## Local Development

Backend:

```bash
cd apps/speclist-api
go run ./cmd/speclist-api
```

Frontend:

```bash
cd apps/speclist-web
npm install
npm run dev
```

## Exported Artifacts

Speclist currently exports reviewed drafts in two formats:

- `openspec-markdown`
  writes a markdown file plus a citation sidecar JSON file
- `fastspec-yaml`
  writes a typed YAML draft document plus the same citation sidecar JSON file

OpenSpec-aware export currently supports writing into active change targets for:

- `proposal.md`
- `design.md`
- `tasks.md`
- `specs/<capability>/spec.md`

When exporting into those OpenSpec targets, Speclist now renders typed artifact templates instead of a single generic markdown draft shape.

The backend requires an explicit target directory and target name for every export.
Edited drafts are normalized during export, and export rejects blank titles or empty section headings/bodies.

## Retrieval Filters

Speclist search and draft generation now support the same compact retrieval filters:

- `kinds`
  limit results to selected source kinds such as `docx`, `confluence`, or `spec`
- `origin`
  choose between `imported` sources and `repository` specs
- `location_contains`
  keep only sources whose location contains the provided substring

The workbench keeps these filters visible and applies them to both retrieval and drafting so the generated draft stays grounded in the intended subset of sources.

## Draft Editing

After generation, the workbench keeps the draft in editable state:

- title and summary can be rewritten before export
- section headings and section bodies can be edited directly
- citations can be adjusted as one entry per line
- reviewers can add or remove sections before export

This editing flow remains in-memory only. The backend validates the edited draft during export so invalid reviewer edits do not produce partial artifacts.

## Draft Presets

Speclist draft generation now supports a small preset set:

- `general`
  the original grounded draft structure with `Why`, `Context`, and `Proposed Requirements`
- `proposal`
  starts with `Why`, `What Changes`, and `Impact`
- `design`
  starts with `Context`, `Goals / Non-Goals`, `Decisions`, and `Risks / Trade-offs`
- `requirements`
  starts with `Why`, `Requirements`, and `Scenarios`

The selected preset is stored on the generated draft payload and remains visible during review and export.

## Draft Diff Review

Speclist now keeps the originally generated draft alongside the live editable draft in the workbench.

The draft diff panel highlights:

- title changes
- summary changes
- preset changes
- section changes by position
- added or removed sections

This keeps the review step explicit before export without introducing server-side draft versioning.

## Review Flags

Speclist now lets reviewers mark each draft section with a lightweight review state:

- `ready`
- `needs-work`
- `blocked`

Each section can also carry an optional review note. The workbench summarizes all non-ready or noted sections in one review flags panel so unresolved concerns are visible before export.

## Export Readiness

Speclist now computes a compact export-readiness summary in the workbench before export.

It reports:

- `ready`
  no blockers or warnings detected
- `warning`
  non-blocking issues such as `needs-work` sections, review notes, missing citations, or an unchanged generated draft
- `blocked`
  blocking issues such as empty required draft fields or sections explicitly marked `blocked`

This readiness check is advisory and client-side. Backend export validation remains the final guardrail.

## Review Reset

Speclist now lets reviewers reset the current draft back to the original generated draft.

Reset clears:

- current local edits
- section review flags
- export result state

The original generated snapshot remains the reset baseline until the reviewer generates a new draft.

## Section Reordering

Speclist now lets reviewers move draft sections up or down during review.

Reordering preserves:

- section content
- section review flags
- the rest of the in-memory review state

This keeps structural review lightweight without introducing drag-and-drop or server-side ordering logic.

## Section Duplication

Speclist now lets reviewers duplicate a draft section during review.

Duplication preserves:

- section content
- section review flags
- the rest of the in-memory review state

The duplicated section is inserted next to the original so the reviewer can refine the copy immediately.

## Citation Inspection

Speclist now lets reviewers inspect grounded source context for any citation visible in:

- retrieval results
- draft section citation lists

The citation inspector resolves the citation string back to the stored source and chunk, then shows:

- source title
- source location
- cited section
- cited excerpt
- source metadata when available

This keeps evidence review inside the workbench instead of forcing reviewers to infer grounding from citation labels alone.

## Source Detail Review

Speclist now lets operators inspect indexed sources directly from the source list.

The source detail panel shows:

- source title and kind
- source location
- source metadata when available
- the full chunk inventory for that source

Each chunk in the source detail panel can also reopen citation inspection, so corpus review and evidence review remain connected.
