## Why

Speclist drafts now support reviewer-side editing, but citations are still just text labels. Reviewers need a direct way to inspect the grounded chunk behind a citation while editing a draft so they can verify evidence before export.

## What Changes

- Add backend citation inspection lookup by citation string.
- Expose a workbench action to inspect citations from search results and draft sections.
- Show the cited source, section, excerpt, and metadata in a focused review panel.

## Capabilities

### New Capabilities
- `speclist-citation-inspector`: inspect the grounded source chunk behind a citation

### Modified Capabilities
- `speclist-workbench`: citations become actionable review elements instead of plain text
- `spec-rag-retrieval`: retrieval evidence can be reopened after draft generation

## Impact

- Affected code: `apps/speclist-api`, `apps/speclist-web`
- Affected behavior: reviewers can verify citation grounding during draft review
- Affected tests: backend service and HTTP coverage for citation lookup
