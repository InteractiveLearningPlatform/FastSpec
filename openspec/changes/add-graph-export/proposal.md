## Why

FastSpec can now parse, inspect, and validate a spec tree, but it still lacks a normalized graph export that downstream agents can use for architecture reasoning or generation workflows. The next useful runtime surface is a stable graph view of projects, modules, workflows, and internal module dependencies.

## What Changes

- Add a `graph` command to the FastSpec CLI.
- Export a normalized graph for a validated FastSpec tree, including project, module, and workflow nodes.
- Represent internal module relationships as explicit edges.
- Support both human-readable and JSON output for graph export.
- Add tests for graph output and invalid-tree behavior.

## Capabilities

### New Capabilities
- `graph-export`: Export a normalized FastSpec project graph for agent and tooling consumption.

### Modified Capabilities
- `validation-findings`: Require graph export to operate on a validation-clean FastSpec tree.
- `json-cli-output`: Extend machine-readable output support to graph export.

## Impact

- Affected code: `apps/fastspec-cli`, `crates/fastspec-core`
- Affected docs: CLI usage docs
- Affected tests: new graph command integration tests
