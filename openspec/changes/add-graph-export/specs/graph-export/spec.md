## ADDED Requirements

### Requirement: CLI MUST export a normalized project graph
The FastSpec CLI MUST provide a `graph` command that exports a normalized graph for a validation-clean FastSpec tree.

#### Scenario: Export graph from valid tree
- **WHEN** a contributor or agent runs `fastspec graph <path>` on a valid FastSpec tree
- **THEN** the CLI exits successfully
- **AND** it emits a graph containing the project, declared modules, workflows, and internal dependency relationships

#### Scenario: Reject graph export from invalid tree
- **WHEN** a contributor or agent runs `fastspec graph <path>` on a FastSpec tree with validation findings
- **THEN** the CLI exits with a non-zero status
- **AND** it reports that graph export requires a validation-clean tree

### Requirement: Graph MUST use explicit nodes and edges
The exported graph MUST represent structure through normalized node and edge records rather than only free-form text.

#### Scenario: Project and module nodes
- **WHEN** the CLI exports a graph from a valid FastSpec project
- **THEN** the graph includes a project node and module nodes for every declared module document

#### Scenario: Workflow nodes
- **WHEN** the project declares workflows
- **THEN** the graph includes workflow nodes and edges connecting them to the project

#### Scenario: Internal dependency edges
- **WHEN** one module depends on another declared module in the same project
- **THEN** the graph includes an explicit dependency edge between those module nodes

### Requirement: Graph MUST support JSON output
The `graph` command MUST support machine-readable JSON output.

#### Scenario: Graph as JSON
- **WHEN** a contributor or agent runs `fastspec graph --json <path>`
- **THEN** the CLI emits valid JSON containing the normalized nodes and edges
