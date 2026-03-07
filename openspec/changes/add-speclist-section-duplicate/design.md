## Context

The current workbench supports section editing, reordering, and reset, but duplicating a section still requires copy-pasting content manually. Since section state already lives entirely in memory, duplication can be implemented as another local section-list transformation.

## Goals / Non-Goals

**Goals:**
- Duplicate a section with one action.
- Insert the copy adjacent to the original section.
- Copy section review metadata alongside the duplicated section.

**Non-Goals:**
- Add named templates for duplication.
- Duplicate whole drafts.
- Introduce backend support for duplication.

## Decisions

Insert duplicates immediately after the original section.
Rationale: this keeps the copy visually close to the source section so the reviewer can refine it right away.

Copy review flags to the duplicated section.
Rationale: the duplicated section starts as the same review object and can be adjusted afterward.

Keep duplicated content identical to the source section.
Rationale: the point is to branch from the current section state, not to infer a new variation automatically.

## Risks / Trade-offs

[Copied review flags may not always be desirable] -> Start by preserving them for consistency and let reviewers change them immediately if needed.

[Repeated duplication may create noisy drafts] -> The existing reset and remove controls remain the cleanup path.
