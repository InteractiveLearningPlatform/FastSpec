## Why

`project.spec.workflows` declares workflows as bare `IdPurpose` pairs — an id and a single purpose line. Every other spec kind (`ModuleSpec`, `AgentCapabilitySpec`) has a corresponding document with a rich body. Workflows are a gap: they appear as nodes in the graph and steps in the plan but carry no retrievable structured detail. An agent querying a workflow spec gets nothing beyond the purpose line from the project declaration.

## What Changes

- Add `WorkflowSpecBody` and `WorkflowSpecDocument` to `fastspec-model`
- Wire `WorkflowSpec` through the full pipeline in `fastspec-core`:
  - `validate_findings`: `missing_workflow_document` and `undeclared_workflow_document` cross-document rules (symmetric with module and capability rules)
  - `export_graph`: build workflow nodes from actual `WorkflowSpec` documents instead of inline `IdPurpose` data; use workflow document title in node
  - `export_plan`: workflow steps reflect document-level data
  - `generate_scaffold`: write `workflows/<id>/README.md` from the workflow document body instead of the project's `IdPurpose.purpose`
- Add `WorkflowSpec` kind to the CLI's `summary`, `inspect`, `graph`, `plan`, `generate` commands (all flow through the pipeline, so no CLI changes needed beyond the model/core updates)
- Add `workflow.fastspec.yaml` template to `templates/`
- Add `workflow:plan.fastspec.yaml` and `workflow:generate.fastspec.yaml` to the archlint-reproduction example
- Update example `project.fastspec.yaml` to declare the workflow docs (already declared by id; now they need matching files)
- Update all tests that count documents or check artifacts

## Capabilities

### New Capabilities
- `workflow-spec-document`: First-class `WorkflowSpec` document kind with structured body (steps, triggers, inputs, outputs); validated, graphed, planned, and scaffolded symmetrically with `ModuleSpec` and `AgentCapabilitySpec`

### Modified Capabilities
- None. Existing pipeline rules and tests are extended, not changed in contract.

## Impact

- `crates/fastspec-model/src/lib.rs`: new structs, new enum variant
- `crates/fastspec-core/src/lib.rs`: two new validation rules, graph/plan/scaffold updates, new render function
- `templates/`: new `workflow.fastspec.yaml`
- `examples/archlint-reproduction/specs/`: two new workflow spec files
- `apps/fastspec-cli/tests/`: count updates, new workflow assertions
