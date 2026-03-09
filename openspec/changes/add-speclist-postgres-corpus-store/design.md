## Context

Speclist currently uses `FileStore`, which serializes the full corpus into one JSON file. That keeps the first slice simple, but it does not match the PostgreSQL operational baseline that now exists in the repository's compose and Helm scaffolding. The domain already exposes a `CorpusStore` interface, so this change can add a PostgreSQL adapter without forcing service-layer rewrites.

This slice stays intentionally narrow: it only replaces corpus persistence for `SourceDocument` and `Chunk`. Marketplace, analytics, queue, and retrieval-index adapters remain later changes.

## Goals / Non-Goals

**Goals:**
- Add a PostgreSQL-backed `CorpusStore` implementation behind the existing domain interface.
- Keep `FileStore` available so local development still works without PostgreSQL.
- Select the store at runtime through explicit environment configuration.
- Persist source documents and chunks in a queryable schema while preserving current read/write behavior.

**Non-Goals:**
- Implement PostgreSQL adapters for marketplace catalog, analytics, cache, queue, or search assets.
- Introduce a full migration framework.
- Change the Speclist HTTP API surface or retrieval semantics.

## Decisions

Use `database/sql` with the `pgx` stdlib driver rather than a heavier ORM.
Rationale: the current service is small, and the store contract is simple. `database/sql` plus `pgx` is enough to manage schema bootstrap and straightforward CRUD without introducing extra abstractions.
Alternative considered: `gorm`. Rejected because it adds more framework surface than this slice needs.

Store document metadata and chunk metadata as JSONB columns.
Rationale: the current domain already models metadata as `map[string]string`, so JSONB preserves flexibility without premature schema expansion.
Alternative considered: normalized metadata tables. Rejected because it is unnecessary for the first persistence slice.

Bootstrap the schema on store initialization.
Rationale: the platform ops baseline now includes local compose and CI validation, so the adapter should be able to create its own tables in empty databases.
Alternative considered: manual SQL setup. Rejected because it adds operator friction before the storage baseline is proven.

Keep runtime selection explicit with `SPECLIST_STORE_KIND=file|postgres` and `SPECLIST_POSTGRES_DSN`.
Rationale: the file store remains useful for lightweight local work, but PostgreSQL must be opt-in rather than silently selected from partial config.
Alternative considered: infer PostgreSQL mode from DSN presence only. Rejected because explicit store selection is easier to reason about in CI and compose environments.

## Risks / Trade-offs

[Schema bootstrap drifts from future migration tooling] -> Keep the bootstrap SQL minimal and isolated inside the adapter so it can be replaced later.

[PostgreSQL tests become flaky if they need a real database] -> Use SQL mock coverage for adapter behavior in this slice and defer live database tests to compose-backed integration coverage.

[Document replacement semantics could diverge from the file store] -> Preserve the current `Save` behavior by upserting documents on ID and replacing all stored chunks for that document in one transaction.

## Migration Plan

1. Add the PostgreSQL adapter and store selection in the API bootstrap path.
2. Default existing environments to `file` so the change is non-breaking.
3. Enable PostgreSQL mode in compose or CI by setting `SPECLIST_STORE_KIND=postgres` and a valid DSN.
4. Roll back by switching `SPECLIST_STORE_KIND` back to `file`.

## Open Questions

- Should the next slice add a one-time file-to-PostgreSQL import helper for existing local corpus JSON files?
- Should chunk text eventually be split into a retrieval-oriented table separate from authoring/source storage?
