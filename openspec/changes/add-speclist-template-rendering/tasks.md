## 1. Backend rendering

- [x] 1.1 Add template-aware renderers for proposal, design, tasks, and spec OpenSpec targets
- [x] 1.2 Route OpenSpec-targeted export through the correct renderer based on artifact type
- [x] 1.3 Keep generic filesystem export behavior intact for non-targeted markdown output

## 2. Verification

- [x] 2.1 Add backend tests covering proposal, tasks, and spec template rendering
- [x] 2.2 Verify Go tests still pass after the renderer split

## 3. Docs and UI

- [x] 3.1 Update Speclist docs to describe typed OpenSpec artifact rendering
- [x] 3.2 Verify the frontend build still passes without UI regressions
