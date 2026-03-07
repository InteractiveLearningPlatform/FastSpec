## Why

FastSpec can now export a normalized project graph, but agents still need to infer an execution order from that structure before they can act on it. The next useful runtime capability is a deterministic plan export that turns a validated graph into an ordered implementation sequence.

## What Changes

- Add a `plan` command to the FastSpec CLI.
- Derive an ordered implementation plan from a validation-clean FastSpec graph.
- Include phase and dependency information for project setup, modules, and workflows.
- Support both human-readable and JSON output for plan export.
- Add tests for successful plan generation and invalid-tree rejection.

## Capabilities

### New Capabilities
- `plan-export`: Export an ordered implementation plan derived from a validated FastSpec graph.

### Modified Capabilities
- `graph-export`: Reuse the normalized graph as the source for planning.
- `json-cli-output`: Extend machine-readable output support to plan export.

## Impact

- Affected code: `apps/fastspec-cli`, `crates/fastspec-core`
- Affected docs: CLI usage docs
- Affected tests: new plan command integration tests
