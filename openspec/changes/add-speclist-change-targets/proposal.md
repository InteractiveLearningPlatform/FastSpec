## Why

Speclist can now export reviewed drafts into durable files, but it still treats export targets as generic filesystem paths. The next useful step is to let operators export directly into active OpenSpec change artifacts so reviewed drafts can land in `proposal.md`, `design.md`, `tasks.md`, or `specs/<capability>/spec.md` without manual path construction.

## What Changes

- Add repo-aware OpenSpec change target support to Speclist export.
- List active OpenSpec changes through the backend for the workbench.
- Support exporting reviewed drafts directly into proposal, design, tasks, and spec artifact paths.
- Extend the workbench export UI with an OpenSpec change target mode.

## Capabilities

### New Capabilities
- `speclist-change-targets`: export reviewed drafts directly into active OpenSpec change artifact paths

### Modified Capabilities
- `draft-export`: export now supports OpenSpec-aware artifact targets in addition to generic filesystem destinations
- `speclist-workbench`: the export UI now offers repo-aware OpenSpec change targeting

## Impact

- Affected code: `apps/speclist-api`, `apps/speclist-web`
- Affected workflow: reviewed Speclist drafts can now land directly in active OpenSpec changes
- Affected docs: Speclist workflow and export behavior documentation
