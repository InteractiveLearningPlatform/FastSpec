## 1. Model layer (fastspec-model)

- [x] 1.1 Add `WorkflowSpecBody` struct with `purpose: String`, `steps: Vec<NamedItem>`, `inputs: Vec<NamedItem>`, `outputs: Vec<NamedItem>`, `triggers: Vec<String>` (all vec fields default empty)
- [x] 1.2 Add `WorkflowSpecDocument` struct (apiVersion, kind, metadata, spec)
- [x] 1.3 Add `FastSpecDocument::Workflow(WorkflowSpecDocument)` variant; implement `kind()` and `metadata()` arms; add `spec_detail_lines` for Workflow
- [x] 1.4 Add `SpecKind::Workflow` variant with `as_str() = "WorkflowSpec"`
- [x] 1.5 Add `"WorkflowSpec"` branch in `parse_document` match

## 2. Core pipeline (fastspec-core)

- [x] 2.1 Add two validation rules in `validate_findings`: `missing_workflow_document` (project declares workflow but no doc exists) and `undeclared_workflow_document` (doc exists but not declared in project)
- [x] 2.2 In `export_graph`: pre-collect workflow document paths and titles; build workflow nodes from document metadata (title from `WorkflowSpecDocument.metadata.title`); add `DefinesWorkflow` edges from project to workflow nodes
- [x] 2.3 In `generate_scaffold`: collect workflow documents; change scaffold output from `workflows/<id>.md` to `workflows/<id>/README.md`; write content from `WorkflowSpec` body via new `render_workflow_readme` function
- [x] 2.4 Add `render_workflow_readme` function that renders purpose, steps, inputs, outputs, triggers sections

## 3. Examples and templates

- [x] 3.1 Add `examples/archlint-reproduction/specs/workflows/plan.fastspec.yaml` with realistic content for the "plan" workflow
- [x] 3.2 Add `examples/archlint-reproduction/specs/workflows/generate.fastspec.yaml` with realistic content for the "generate" workflow
- [x] 3.3 Add `templates/workflow.fastspec.yaml` showing all fields with example values

## 4. Tests

- [x] 4.1 Update `validates_archlint_example_tree` and `reports_clean_validation_for_example_tree` in fastspec-core (document count: 4 → 6)
- [x] 4.2 Update `exports_graph_for_example_tree` to assert workflow nodes come from documents
- [x] 4.3 Update `generates_scaffold_for_example_tree` to use new path `workflows/plan/README.md` instead of `workflows/plan.md`
- [x] 4.4 Update `json_output.rs` summary test (3→4 was already done; now 4→6)
- [x] 4.5 Update `generate_output.rs` E2E test path assertion (`workflows/generate.md` → `workflows/generate/README.md`)
- [x] 4.6 Add unit test for `missing_workflow_document` and `undeclared_workflow_document` validation rules
