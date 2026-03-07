# speclist-api

Go backend microservice for Speclist.

Responsibilities:

- ingest DOCX files into a spec-ready corpus
- import Confluence pages into the same corpus
- index existing FastSpec and OpenSpec artifacts from the repo
- expose retrieval bundles and draft-generation endpoints
- export reviewed drafts into OpenSpec markdown or FastSpec YAML files

Architecture:

- hexagonal structure with application services in `internal/app`
- domain types in `internal/domain`
- adapters for HTTP, storage, document import, and spec indexing in `internal/adapters`

Run locally:

```bash
go run ./cmd/speclist-api
```

Environment:

- `SPECLIST_ADDR` default: `:8080`
- `SPECLIST_DATA_FILE` default: `./data/corpus.json`
- `SPECLIST_REPO_ROOT` default: `../..`
