# Speclist

Speclist is the ingestion and retrieval product surface in this repo.

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
5. Export accepted drafts into durable OpenSpec markdown or FastSpec YAML files.
6. Optionally target an active OpenSpec change artifact directly instead of a generic output path.

## Design Constraints

- Keep retrieval compact and interactive.
- Preserve provenance for every returned excerpt.
- Keep the backend domain independent from storage, HTTP, and source-specific adapters.
- Treat generated drafts as reviewable candidates, not final truth.
- Require explicit export destinations and avoid silent overwrite of existing files.

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
  writes a YAML draft file plus the same citation sidecar JSON file

OpenSpec-aware export currently supports writing into active change targets for:

- `proposal.md`
- `design.md`
- `tasks.md`
- `specs/<capability>/spec.md`

The backend requires an explicit target directory and target name for every export.
