## ADDED Requirements

### Requirement: CLI MUST expose validation findings
The FastSpec CLI MUST provide a `validate` command that evaluates one FastSpec file or tree and returns explicit findings rather than only parse success or failure.

#### Scenario: Validate clean tree
- **WHEN** a contributor or agent runs `fastspec validate <path>` on a valid FastSpec tree
- **THEN** the CLI exits successfully
- **AND** it reports that no validation findings were produced

#### Scenario: Validate invalid tree
- **WHEN** a contributor or agent runs `fastspec validate <path>` on a FastSpec tree with semantic problems
- **THEN** the CLI exits with a non-zero status
- **AND** it reports explicit findings describing the violations

### Requirement: Validation MUST detect core cross-document consistency issues
The validation layer MUST detect at least duplicate document identifiers, mismatches between project-declared modules and module documents, and internal module dependency references to undeclared project modules.

#### Scenario: Duplicate identifier
- **WHEN** two supported FastSpec documents in the same validation scope share the same identifier
- **THEN** validation reports a duplicate-identifier finding

#### Scenario: Missing module document
- **WHEN** a project declares a module identifier without a matching module document in the tree
- **THEN** validation reports a missing-module-document finding

#### Scenario: Undeclared module document
- **WHEN** a module document exists but its identifier is not declared in the project module list
- **THEN** validation reports an undeclared-module-document finding

#### Scenario: Internal dependency target missing from project modules
- **WHEN** a module dependency references an identifier intended to be another project module but that identifier is not declared in the project module list
- **THEN** validation reports an invalid-module-dependency finding

### Requirement: Validation MUST support JSON output
The `validate` command MUST support machine-readable JSON output for findings.

#### Scenario: Validation findings as JSON
- **WHEN** a contributor or agent runs `fastspec validate --json <path>`
- **THEN** the CLI emits valid JSON containing the findings and overall validity result

