## Context

The workbench already has a sequential outline navigation function. Keyboard shortcuts should reuse that function rather than creating a second navigation path. The only significant safeguard is preventing shortcuts from firing while the user is typing into form controls.

## Goals / Non-Goals

**Goals:**
- Trigger previous and next outline navigation from the keyboard.
- Keep the shortcut behavior inactive while typing in inputs, textareas, or selects.
- Reuse the existing sequential navigation flow.

**Non-Goals:**
- Add a configurable shortcut map.
- Add shortcuts for every review action.
- Introduce global hotkey infrastructure.

## Decisions

Use document-level key handling that delegates to the existing outline navigation helper.
Rationale: the review surface is already a single-page client UI, and the shortcut target is global within the active draft context.

Use `[` for previous and `]` for next.
Rationale: these keys are unlikely to conflict with text editing conventions for this app and map cleanly to backward/forward stepping.

Ignore shortcut handling when the event target is an input, textarea, or select.
Rationale: reviewers should be able to type normally without accidental navigation.

Only enable shortcuts when a draft exists and the outline has visible entries.
Rationale: shortcuts should do nothing outside the active review flow.

## Risks / Trade-offs

[Document-level shortcuts can conflict with future global actions] -> Reusing a narrow key pair and centralizing the guardrails keeps the scope contained.

[Some keyboards make bracket keys less convenient] -> The feature remains additive because button-based navigation still exists.
