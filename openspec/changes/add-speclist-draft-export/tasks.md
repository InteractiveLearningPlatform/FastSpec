## 1. Backend export contract

- [x] 1.1 Add draft export request, response, and artifact types to the Speclist domain
- [x] 1.2 Implement OpenSpec markdown and FastSpec YAML renderers plus citation sidecar output
- [x] 1.3 Add filesystem export logic with explicit target directory handling and overwrite protection

## 2. API and verification

- [x] 2.1 Expose a draft export endpoint in the HTTP adapter
- [x] 2.2 Add backend tests covering successful export and overwrite rejection
- [x] 2.3 Verify the backend suite still passes with the new export flow

## 3. Workbench and docs

- [x] 3.1 Add frontend export controls and export result rendering to the reviewed draft view
- [x] 3.2 Update Speclist docs to describe exporting reviewed drafts into durable artifacts
- [x] 3.3 Verify the frontend production build still passes
