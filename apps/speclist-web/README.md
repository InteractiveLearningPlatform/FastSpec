# speclist-web

React workbench for Speclist.

Capabilities:

- upload DOCX source documents
- import Confluence pages
- index local FastSpec/OpenSpec artifacts
- search grounded retrieval context
- review generated draft specs and their citations
- export reviewed drafts into durable files through the backend
- target active OpenSpec change artifacts during export

Run locally:

```bash
npm install
npm run dev
```

The frontend expects the API at `http://localhost:8080` by default. Override with `VITE_API_BASE`.
