## ADDED Requirements

### Requirement: CLI MUST generate a starter scaffold
The FastSpec CLI MUST provide a `generate` command that writes a deterministic starter scaffold for a validation-clean FastSpec tree.

#### Scenario: Generate scaffold
- **WHEN** a contributor or agent runs `fastspec generate --out <dir> <path>` on a valid FastSpec tree
- **THEN** the CLI exits successfully
- **AND** it writes a starter scaffold to the requested output directory

#### Scenario: Reject invalid tree generation
- **WHEN** a contributor or agent runs `fastspec generate --out <dir> <path>` on a FastSpec tree with validation findings
- **THEN** the CLI exits with a non-zero status
- **AND** it reports that generation requires a validation-clean tree

### Requirement: Generated scaffold MUST include core project structure
The generated scaffold MUST include a project-level file, per-module directories with stub content, workflow stub files, and a machine-readable manifest of generated artifacts.

#### Scenario: Module scaffold exists
- **WHEN** the project defines modules
- **THEN** the generated output includes one stub directory per module

#### Scenario: Workflow scaffold exists
- **WHEN** the project defines workflows
- **THEN** the generated output includes one stub file per workflow

### Requirement: Generation MUST support JSON reporting
The `generate` command MUST support machine-readable JSON output describing what was written.

#### Scenario: Generate scaffold with JSON report
- **WHEN** a contributor or agent runs `fastspec generate --json --out <dir> <path>`
- **THEN** the CLI emits valid JSON describing the generated files and output directory
