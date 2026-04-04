## Context

`ModuleSpec` and `AgentCapabilitySpec` are both first-class document kinds: they have structured YAML bodies, are validated against project declarations, appear as graph nodes built from their own metadata, and produce scaffold artifacts from their spec content. `WorkflowSpec` is not — workflow nodes in the graph come from project-level `IdPurpose` data, scaffold stubs use the single `purpose` string from the project, and there is no cross-document validation. This is the last asymmetry in the pipeline.

## Goals / Non-Goals

**Goals:**
- Add `WorkflowSpecBody` / `WorkflowSpecDocument` to `fastspec-model` with useful structured fields
- Wire `WorkflowSpec` through the full pipeline: parse → validate → graph → plan → scaffold
- Add two validation rules (`missing_workflow_document`, `undeclared_workflow_document`)
- Build workflow graph nodes from document metadata (title from doc, not from IdPurpose.purpose)
- Scaffold `workflows/<id>/README.md` from document body fields
- Add example workflow documents and a `templates/workflow.fastspec.yaml`

**Non-Goals:**
- Workflow execution or runtime semantics
- Cross-workflow dependencies or ordering rules
- Breaking the existing `IdPurpose`-based workflow node shape (the `GraphNodeKind::Workflow` variant stays)

## Decisions

**WorkflowSpecBody fields**: `purpose` (required string), `steps` (Vec<NamedItem>), `inputs` (Vec<NamedItem>), `outputs` (Vec<NamedItem>), `triggers` (Vec<String>). Reuses the existing `NamedItem` type from `fastspec-model` (name + description). This mirrors ModuleSpec's `inputs`/`outputs` pattern and keeps the model coherent without new types.

**Graph node title**: Switch from `workflow.purpose` (from project IdPurpose) to `workflow_document.metadata.title`. This requires pre-collecting workflow document paths and titles before the main graph-building loop — the same pattern used for modules and capabilities.

**Validation gating**: `export_graph` already requires a validation-clean tree. Adding workflow doc rules to `validate_findings` automatically gates graph/plan/scaffold. No extra gating logic needed.

**Scaffold output path change**: Currently `generate_scaffold` writes `workflows/<id>.md` (a flat file). With WorkflowSpec documents, it will write `workflows/<id>/README.md` (a directory per workflow, matching the module pattern). This is a **behavioral change** for existing trees that have workflows — existing tests will need path updates.

**No WorkflowSpec in PlanPhase**: Workflow steps already exist as `PlanPhase::Workflow`. No new phase variant is needed; the existing variant is correct.

## Risks / Trade-offs

**Test churn**: Every test that counts documents, artifact paths, or checks "workflows/plan.md" will break. The count goes from 4 → 6 documents (adding 2 workflow docs to archlint-reproduction), and the path changes from `workflows/plan.md` to `workflows/plan/README.md`. These are mechanical updates.

**Backwards compatibility**: Any user currently relying on `workflows/<id>.md` scaffold output will see the path change to `workflows/<id>/README.md`. This is acceptable — the codebase has no external callers and the change is the right long-term shape.
