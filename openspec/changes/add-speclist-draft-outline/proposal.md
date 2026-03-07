## Why

Speclist reviewers can now collapse sections individually or all at once, but long drafts still require too much scrolling to find the right section. A compact outline is the next useful navigation aid because it lets reviewers jump directly to the part of the draft they want to inspect.

## What Changes

- Add a draft outline panel that lists the current sections in order.
- Let reviewers jump from the outline to a specific section in the draft.
- Reopen collapsed sections when navigation targets them.

## Capabilities

### New Capabilities
- `speclist-draft-outline`: navigate the current draft through a compact section index

### Modified Capabilities
- `speclist-workbench`: reviewers can navigate long drafts without manual scrolling
- `speclist-section-collapse`: collapsed sections can be reopened through outline navigation

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: long draft review becomes faster and more navigable
- Affected tests: frontend build verification
