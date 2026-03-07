## Why

Speclist is currently a focused ingestion and export tool, but the product direction you described is much broader: a production platform and marketplace for specs with strong ML, RAG, and NLP pipelines. That requires a dedicated future change covering scalable search/indexing, production security gates, and deployment architecture rather than piecemeal feature additions.

## What Changes

- Evolve Speclist toward a marketplace/platform for specs instead of a single-workbench tool.
- Add ML, RAG, and NLP pipeline requirements for spec ingestion, normalization, retrieval, ranking, and generation.
- Define a DB-agnostic architecture that can use PostgreSQL, ClickHouse, and Valkey together.
- Add a scalable vector and hybrid retrieval layer for specs, code, GitHub markdown, and language-oriented intermediate representations.
- Define production-grade Docker Compose deployment with Traefik, CrowdSec, and Trivy-based security validation gates.
- Require secret handling and configuration validation so insecure passwords or hardcoded sensitive values block startup.
- Add CI/CD requirements for GitHub and Kubernetes deployment using Helm and Vault-backed secret management.

## Capabilities

### New Capabilities
- `speclist-platform-marketplace`: marketplace and platform behaviors for publishing, indexing, discovering, and reusing specs
- `spec-ml-pipelines`: ML, RAG, and NLP pipelines for ingestion, enrichment, ranking, and generation
- `spec-search-infrastructure`: scalable indexing and retrieval architecture across structured specs, markdown, code, and IR-oriented artifacts
- `production-security-gates`: production deployment, secret validation, and startup-blocking security checks
- `platform-delivery-pipeline`: CI/CD and Kubernetes delivery for the Speclist platform

### Modified Capabilities
- `spec-rag-retrieval`: retrieval expands from simple document/spec lookup into scalable hybrid search with code and IR-aware indexing
- `speclist-workbench`: the workbench becomes one surface of a broader Speclist platform rather than the whole product

## Impact

- Affected architecture: backend services, indexing/storage topology, security controls, deployment stack
- Affected infrastructure: PostgreSQL, ClickHouse, Valkey, vector retrieval layer, Docker Compose, Kubernetes, Helm, Vault
- Affected operations: secret management, config linting, startup validation, CI/CD
- Affected docs/specs: platform architecture, search stack selection, production deployment, security policy
