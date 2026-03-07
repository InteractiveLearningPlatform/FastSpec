## 1. Typed parsing model

- [x] 1.1 Add YAML parsing dependencies and define typed Rust models for supported FastSpec documents in `fastspec-model`
- [x] 1.2 Implement parsing helpers that load a FastSpec YAML file into the correct typed document variant

## 2. Validation and inspection

- [x] 2.1 Update `fastspec-core` to validate parsed documents and report malformed required fields
- [x] 2.2 Extend the CLI with an inspection command that renders typed document details for a file or directory

## 3. Verification and examples

- [x] 3.1 Add tests for successful parsing of the example specs and failure cases for invalid documents
- [x] 3.2 Update CLI or example docs so contributors can use the new typed inspection workflow
