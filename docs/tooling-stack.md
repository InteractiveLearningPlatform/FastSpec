---
area: docs
component: tooling-stack
kind: reference
status: active
language: multiple
framework: openspec
tags:
  - openspec
  - docker-compose
  - traefik
  - crowdsec
  - trivy
---

# Tooling Stack Context

This document stores durable repo-local context for the workflow and infrastructure tools FastSpec is likely to reference.

## OpenSpec

OpenSpec is the default workflow layer for this repo.

- It is explicitly positioned as a lightweight spec framework with an artifact-guided flow: propose, apply, archive.
- Its philosophy is fluid rather than rigid, iterative rather than waterfall, and designed for brownfield as well as greenfield work.
- Each change gets its own folder with proposal, design, tasks, and scoped specs, which makes it a good fit for AI-assisted implementation.
- It requires Node.js 20.19.0 or higher and supports many AI tools through generated prompts and skills.

For FastSpec, OpenSpec is not the durable system model. It is the execution workflow used to create and evolve that model.

## Docker Compose

Docker Compose is the preferred local orchestration format for multi-service examples.

- Docker describes Compose as a tool for defining and running multi-container applications.
- The key advantage is keeping services, networks, and volumes together in a single YAML file.
- Compose is suitable across development, testing, CI, and production-oriented workflows.
- It supports lifecycle commands for starting, stopping, rebuilding, viewing status, streaming logs, and running one-off commands on services.

For FastSpec, Compose is a good target shape when an example system needs a concise, multi-container deployment model.

## Traefik

Traefik is the preferred reverse proxy and routing reference for containerized examples.

- Traefik describes itself as an application proxy that automatically discovers infrastructure configuration and updates routes dynamically.
- Core concepts are entrypoints, routers, middleware, services, and providers.
- Providers are how Traefik discovers route information; Docker and file providers are especially relevant for local Compose stacks.
- Entrypoints define the listening ports and protocols.
- Routers connect requests to services and may use middleware before forwarding.

For FastSpec examples, Traefik is valuable when the spec needs an explicit edge layer that reacts to container metadata instead of static hand wiring.

## CrowdSec

CrowdSec is the preferred detection and remediation reference for security-oriented deployment examples.

- CrowdSec's Security Engine analyzes logs for behavior-based scenarios and supports remediation and AppSec protection.
- In Docker deployments, CrowdSec requires persisted data under `/var/lib/crowdsec/data`, and the docs recommend Compose for most production use cases.
- Containerized deployments need access to log sources through mounted volumes.
- The Local API and remediation components are separate concerns; detection alone is not enough if you also want blocking behavior.
- CrowdSec's Traefik guidance emphasizes correct client IP forwarding and trusted proxy configuration so remediation acts on real source IPs.

For FastSpec, CrowdSec is useful as a model for describing detection pipelines, shared volumes, log acquisition, and trust boundaries between reverse proxy and remediation layers.

## Trivy

Trivy is the preferred baseline scanner for code, configuration, and images.

- Aqua describes Trivy as a comprehensive security scanner with multiple targets and scanners.
- Supported targets include container images, filesystems, git repositories, VM images, and Kubernetes.
- Trivy covers vulnerabilities, SBOM/package inventory, misconfigurations, secrets, and licenses.
- It works well as both a local development check and a CI gate.

For FastSpec, Trivy is a useful reference for security validation workflows because it spans repo files, IaC, and container artifacts without forcing separate tools for each surface.

## How These Pieces Fit Together

In a typical local platform example:

1. Compose defines services, networks, and volumes.
2. Traefik fronts the services and handles routing.
3. CrowdSec consumes logs and applies remediation through a bouncer or plugin path.
4. Trivy scans the repo, config, and images before or during CI.
5. OpenSpec remains the workflow layer that plans and evolves the stack.
