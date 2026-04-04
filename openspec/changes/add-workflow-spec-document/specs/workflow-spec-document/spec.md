## ADDED Requirements

### Requirement: WorkflowSpec document kind
`WorkflowSpec` is a first-class FastSpec document kind, symmetric with `ModuleSpec` and `AgentCapabilitySpec`. Each workflow declared in a project's `spec.workflows` array MUST have a corresponding `WorkflowSpec` document reachable from the spec tree.

#### Scenario: parse WorkflowSpec document
- **WHEN** a `.yaml` file with `kind: WorkflowSpec` is parsed
- **THEN** the parser produces a `FastSpecDocument::Workflow(WorkflowSpecDocument)` variant with structured body fields

#### Scenario: WorkflowSpecBody fields
- **WHEN** a `WorkflowSpec` document is parsed
- **THEN** it contains: `purpose` (string, required), `steps` (list of named step items, default empty), `inputs` (list of named items, default empty), `outputs` (list of named items, default empty), `triggers` (list of strings, default empty)

### Requirement: Validation — missing and undeclared workflow documents
The validator enforces a declared-vs-present contract for workflow documents, matching the existing module and capability rules.

#### Scenario: project declares workflow but no document exists
- **WHEN** `project.spec.workflows` contains an id with no matching `WorkflowSpec` document in the tree
- **THEN** `validate_findings` emits a `missing_workflow_document` finding with `Error` severity

#### Scenario: workflow document exists but not declared in project
- **WHEN** a `WorkflowSpec` document exists in the spec tree but its id is not in `project.spec.workflows`
- **THEN** `validate_findings` emits an `undeclared_workflow_document` finding with `Error` severity

#### Scenario: clean tree passes validation
- **WHEN** every declared workflow has a matching document and every workflow document is declared
- **THEN** no workflow-related findings are emitted

### Requirement: Graph — workflow nodes from documents
Workflow nodes in the graph are built from `WorkflowSpec` documents, not from the project's `IdPurpose` data.

#### Scenario: workflow node title from document
- **WHEN** `export_graph` builds workflow nodes
- **THEN** each node's `title` is taken from the `WorkflowSpec` document's `metadata.title`, not from `IdPurpose.purpose`

### Requirement: Scaffold — workflow directory and README
`generate_scaffold` writes a structured `README.md` for each workflow based on the `WorkflowSpec` body.

#### Scenario: scaffold writes workflow README
- **WHEN** `generate_scaffold` runs on a valid tree with workflow documents
- **THEN** `workflows/<id>/README.md` is written with sections for purpose, steps, inputs, outputs, and triggers

### Requirement: Example and template coverage
The archlint-reproduction example and the `templates/` directory are updated.

#### Scenario: example tree includes workflow documents
- **WHEN** the archlint-reproduction spec tree is validated
- **THEN** `plan.fastspec.yaml` and `generate.fastspec.yaml` exist in `specs/workflows/` and pass validation

#### Scenario: workflow template exists
- **WHEN** a user inspects `templates/workflow.fastspec.yaml`
- **THEN** it contains all required and optional fields with example values
