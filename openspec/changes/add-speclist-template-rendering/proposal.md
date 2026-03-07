## Why

Speclist can now export directly into active OpenSpec artifacts, but the rendered content is still generic draft prose. The next useful step is to render typed artifact templates so proposal, design, tasks, and spec exports land closer to the structure that OpenSpec contributors actually expect to refine.

## What Changes

- Add typed template rendering for OpenSpec artifact exports.
- Render proposal, design, tasks, and spec exports differently based on the selected target artifact.
- Preserve citations while shaping output into the expected artifact structure.
- Keep generic filesystem export available, but improve the OpenSpec-targeted rendering path.

## Capabilities

### New Capabilities
- `speclist-template-rendering`: render reviewed drafts into artifact-aware OpenSpec templates

### Modified Capabilities
- `draft-export`: export now renders target-specific OpenSpec artifact content rather than generic draft markdown
- `speclist-change-targets`: OpenSpec change export now uses typed artifact templates for the supported targets

## Impact

- Affected code: `apps/speclist-api`, `apps/speclist-web`
- Affected behavior: exported OpenSpec artifacts become more structured and closer to final workflow shape
- Affected tests: export rendering coverage for each supported artifact class
