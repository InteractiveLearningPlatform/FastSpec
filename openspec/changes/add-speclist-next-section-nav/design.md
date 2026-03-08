## Context

The workbench already has everything needed for sequential navigation: an active section index, a filtered outline list, and a focus helper that expands and scrolls to sections. Guided stepping should reuse that model rather than introducing separate review queues or server state.

## Goals / Non-Goals

**Goals:**
- Let reviewers move to the next outline entry from the current active section.
- Let reviewers move backward through the outline as well.
- Make stepping follow the current filtered outline view.

**Non-Goals:**
- Persist navigation history.
- Add keyboard shortcuts in this slice.
- Create separate queues for blocked or needs-work sections.

## Decisions

Place next and previous controls in the draft outline panel.
Rationale: the behavior is about moving through the outline, not editing section content directly.

Base sequential navigation on the current filtered outline entries.
Rationale: if the reviewer narrows the outline, guided navigation should honor that working set.

When the active section is not present in the filtered outline, start navigation from the first matching entry.
Rationale: filtering can intentionally hide the current active section, so the next valid target should still be reachable in one click.

Wrap navigation at the ends of the filtered list.
Rationale: continuous stepping is more useful during review than hard stops.

## Risks / Trade-offs

[Wrapping can surprise some users] -> The controls stay small and local to the outline, and wrap-around keeps review passes efficient.

[Filtered navigation can feel inconsistent if the active section disappears] -> Starting from the first visible match keeps the behavior deterministic.
