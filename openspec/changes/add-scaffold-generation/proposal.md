## Why

FastSpec can now derive an ordered implementation plan, but it still stops short of producing any concrete output from that plan. The next useful runtime step is to generate a deterministic starter scaffold so agents can move from planning to an initial workspace layout without inventing structure ad hoc.

## What Changes

- Add a `generate` command to the FastSpec CLI.
- Generate a starter scaffold from a validation-clean FastSpec tree and derived plan.
- Write a deterministic directory layout with project, module, and workflow stub files.
- Require an explicit output directory via `--out`.
- Support both human-readable and JSON reporting of generated artifacts.

## Capabilities

### New Capabilities
- `scaffold-generation`: Generate a deterministic starter scaffold from a validated FastSpec plan.

### Modified Capabilities
- `plan-export`: Reuse the ordered plan as the basis for filesystem generation.
- `json-cli-output`: Extend machine-readable output support to scaffold generation.

## Impact

- Affected code: `apps/fastspec-cli`, `crates/fastspec-core`
- Affected docs: CLI usage docs
- Affected tests: scaffold generation integration tests
