## Why

Speclist still persists all imported and indexed source documents to a local JSON file, which blocks multi-instance use and does not align with the PostgreSQL baseline defined in platform ops. The next storage slice should move corpus persistence onto PostgreSQL while keeping the domain contract stable.

## What Changes

- Add a PostgreSQL-backed implementation of the existing `CorpusStore` interface for Speclist source documents and chunks.
- Add runtime configuration to select PostgreSQL persistence via DSN while keeping the file store as a fallback.
- Bootstrap the required PostgreSQL schema automatically so local compose and CI environments can start without a separate migration tool.
- Add focused tests for PostgreSQL persistence behavior and runtime store selection.

## Capabilities

### New Capabilities
- `speclist-postgres-corpus-store`: Persist Speclist corpus documents and chunks in PostgreSQL through the existing store interface.

### Modified Capabilities
- `repository-foundation`: Add a durable runtime persistence option for Speclist without changing the OpenSpec workflow model.

## Impact

- `apps/speclist-api/internal/adapters/storage`
- `apps/speclist-api/cmd/speclist-api`
- `apps/speclist-api/go.mod`
- Speclist runtime environment configuration
