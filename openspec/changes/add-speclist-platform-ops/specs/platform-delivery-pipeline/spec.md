## ADDED Requirements

### Requirement: Platform MUST provide GitHub-driven CI/CD to Kubernetes
The platform MUST provide CI/CD from GitHub into a Kubernetes cluster for production deployments.

#### Scenario: Change delivered from GitHub
- **WHEN** a change is merged through the repository workflow
- **THEN** the CI/CD pipeline can build, validate, and deploy the platform
- **AND** the deployment path targets Kubernetes

### Requirement: Kubernetes delivery MUST use Helm and Vault-backed secret management
The production delivery path MUST use Helm for deployment packaging and Vault-backed secret management for cluster environments.

#### Scenario: Cluster deployment configured
- **WHEN** contributors review the Kubernetes deployment path
- **THEN** they find Helm-based release definitions
- **AND** cluster secrets are managed through Vault-backed flows rather than hardcoded manifests
