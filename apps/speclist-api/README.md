# speclist-api

Go backend microservice for Speclist.

Responsibilities:

- ingest DOCX files into a spec-ready corpus
- import Confluence pages into the same corpus
- index existing FastSpec and OpenSpec artifacts from the repo
- expose retrieval bundles and draft-generation endpoints
- export reviewed drafts into OpenSpec markdown or FastSpec YAML files
- export reviewed drafts directly into active OpenSpec change artifact paths

Architecture:

- hexagonal structure with application services in `internal/app`
- domain types in `internal/domain`
- adapters for HTTP, storage, document import, and spec indexing in `internal/adapters`

Platform roadmap:

- marketplace, hybrid retrieval, and production-ops definitions live in
  `docs/speclist-platform-ops.md`
- the domain package includes marketplace asset and storage port definitions for
  the PostgreSQL, ClickHouse, Valkey, and Qdrant baseline

Run locally:

```bash
go run ./cmd/speclist-api
```

Environment:

- `SPECLIST_STORE_KIND` default: `file`
- `SPECLIST_ADDR` default: `:8080`
- `SPECLIST_DATA_FILE` default: `./data/corpus.json`
- `SPECLIST_POSTGRES_DSN` required when `SPECLIST_STORE_KIND=postgres`
- `SPECLIST_REPO_ROOT` default: `../..`
