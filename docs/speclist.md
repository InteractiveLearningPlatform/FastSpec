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
9. Export accepted drafts into durable OpenSpec markdown or FastSpec YAML files.
10. Optionally target an active OpenSpec change artifact directly instead of a generic output path.
11. Narrow retrieval and drafting with source-kind, source-origin, and location filters when the corpus mixes imported docs and repository specs.

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
