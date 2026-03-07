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
5. Move accepted drafts into durable FastSpec or OpenSpec artifacts.

## Design Constraints

- Keep retrieval compact and interactive.
- Preserve provenance for every returned excerpt.
- Keep the backend domain independent from storage, HTTP, and source-specific adapters.
- Treat generated drafts as reviewable candidates, not final truth.

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
