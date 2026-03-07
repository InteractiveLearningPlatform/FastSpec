import { useEffect, useState } from "react";

const apiBase = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";

const initialConfluence = {
  base_url: "",
  page_id: "",
  token: "",
};

const initialExport = {
  format: "openspec-markdown",
  target_dir: "./exports",
  target_name: "speclist-draft",
};

export default function App() {
  const [sources, setSources] = useState([]);
  const [searchQuery, setSearchQuery] = useState("spec retrieval citations");
  const [searchResults, setSearchResults] = useState([]);
  const [draftTitle, setDraftTitle] = useState("Speclist Draft");
  const [draft, setDraft] = useState(null);
  const [exportResult, setExportResult] = useState(null);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [confluence, setConfluence] = useState(initialConfluence);
  const [exportConfig, setExportConfig] = useState(initialExport);

  useEffect(() => {
    void refreshSources();
  }, []);

  async function refreshSources() {
    try {
      const response = await fetch(`${apiBase}/api/v1/sources`);
      const payload = await response.json();
      setSources(payload.sources ?? []);
    } catch (fetchError) {
      setError(fetchError.message);
    }
  }

  async function handleDocxUpload(event) {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    const formData = new FormData();
    formData.append("file", file);
    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/import/docx`, {
        method: "POST",
        body: formData,
      });
      await assertOk(response);
      setMessage(`Imported ${file.name}`);
      await refreshSources();
    });
    event.target.value = "";
  }

  async function handleConfluenceImport(event) {
    event.preventDefault();
    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/import/confluence`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(confluence),
      });
      await assertOk(response);
      setMessage(`Imported Confluence page ${confluence.page_id}`);
      setConfluence(initialConfluence);
      await refreshSources();
    });
  }

  async function handleSpecIndex() {
    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/index/specs`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({}),
      });
      const payload = await assertOk(response);
      setMessage(`Indexed ${payload.count} repository spec documents`);
      await refreshSources();
    });
  }

  async function handleSearch(event) {
    event.preventDefault();
    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/search`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query: searchQuery, limit: 8 }),
      });
      const payload = await assertOk(response);
      setSearchResults(payload.results ?? []);
      setMessage(`Loaded ${payload.results?.length ?? 0} retrieval results`);
    });
  }

  async function handleDraft(event) {
    event.preventDefault();
    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/drafts`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: searchQuery,
          title: draftTitle,
          format: "openspec-markdown",
          limit: 8,
        }),
      });
      const payload = await assertOk(response);
      setDraft(payload);
      setExportResult(null);
      setExportConfig((current) => ({
        ...current,
        target_name: slugify(payload.title || draftTitle),
      }));
      setMessage(`Generated draft from ${payload.source_count} result(s)`);
    });
  }

  async function handleExport(event) {
    event.preventDefault();
    if (!draft) {
      setError("Generate a draft before exporting.");
      return;
    }

    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/exports`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          draft,
          format: exportConfig.format,
          target_dir: exportConfig.target_dir,
          target_name: exportConfig.target_name,
        }),
      });
      const payload = await assertOk(response);
      setExportResult(payload);
      setMessage(`Exported ${payload.artifacts?.length ?? 0} artifact(s)`);
    });
  }

  async function runAction(action) {
    setLoading(true);
    setError("");
    setMessage("");
    try {
      await action();
    } catch (actionError) {
      setError(actionError.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="page">
      <header className="hero">
        <div>
          <p className="eyebrow">Speclist</p>
          <h1>Spec-focused RAG workbench</h1>
          <p className="lede">
            Ingest DOCX and Confluence docs, index repo specs, retrieve grounded context, and review draft specs with
            citations.
          </p>
        </div>
        <div className="statusCard">
          <span>{loading ? "Working..." : "Idle"}</span>
          <strong>{sources.length} sources</strong>
        </div>
      </header>

      {(message || error) && (
        <div className={error ? "banner error" : "banner success"}>
          {error || message}
        </div>
      )}

      <main className="grid">
        <section className="panel">
          <h2>Import Sources</h2>
          <label className="upload">
            <span>Upload DOCX</span>
            <input type="file" accept=".docx" onChange={handleDocxUpload} />
          </label>

          <form className="stack" onSubmit={handleConfluenceImport}>
            <h3>Import Confluence Page</h3>
            <input
              placeholder="Base URL"
              value={confluence.base_url}
              onChange={(event) => setConfluence({ ...confluence, base_url: event.target.value })}
            />
            <input
              placeholder="Page ID"
              value={confluence.page_id}
              onChange={(event) => setConfluence({ ...confluence, page_id: event.target.value })}
            />
            <input
              placeholder="Token (optional)"
              value={confluence.token}
              onChange={(event) => setConfluence({ ...confluence, token: event.target.value })}
            />
            <button type="submit" disabled={loading}>
              Import page
            </button>
          </form>

          <button className="secondary" type="button" onClick={handleSpecIndex} disabled={loading}>
            Index repository specs
          </button>
        </section>

        <section className="panel">
          <h2>Retrieve Context</h2>
          <form className="stack" onSubmit={handleSearch}>
            <textarea value={searchQuery} onChange={(event) => setSearchQuery(event.target.value)} rows={4} />
            <button type="submit" disabled={loading}>
              Search corpus
            </button>
          </form>
          <div className="results">
            {searchResults.map((result) => (
              <article key={result.chunk.id} className="resultCard">
                <div className="resultMeta">
                  <strong>{result.source.title}</strong>
                  <span>{result.chunk.citation}</span>
                </div>
                <p>{result.chunk.text}</p>
              </article>
            ))}
            {searchResults.length === 0 && <p className="empty">No retrieval bundle loaded yet.</p>}
          </div>
        </section>

        <section className="panel">
          <h2>Create Draft</h2>
          <form className="stack" onSubmit={handleDraft}>
            <input value={draftTitle} onChange={(event) => setDraftTitle(event.target.value)} placeholder="Draft title" />
            <button type="submit" disabled={loading}>
              Draft spec
            </button>
          </form>
          {draft ? (
            <div className="draft">
              <p className="draftSummary">{draft.summary}</p>
              {draft.sections.map((section) => (
                <section key={section.heading} className="draftSection">
                  <h3>{section.heading}</h3>
                  <pre>{section.body}</pre>
                  <ul>
                    {section.citations.map((citation) => (
                      <li key={citation}>{citation}</li>
                    ))}
                  </ul>
                </section>
              ))}
              <form className="stack exportForm" onSubmit={handleExport}>
                <h3>Export Reviewed Draft</h3>
                <select
                  value={exportConfig.format}
                  onChange={(event) => setExportConfig({ ...exportConfig, format: event.target.value })}
                >
                  <option value="openspec-markdown">OpenSpec markdown</option>
                  <option value="fastspec-yaml">FastSpec YAML</option>
                </select>
                <input
                  value={exportConfig.target_dir}
                  onChange={(event) => setExportConfig({ ...exportConfig, target_dir: event.target.value })}
                  placeholder="Target directory"
                />
                <input
                  value={exportConfig.target_name}
                  onChange={(event) => setExportConfig({ ...exportConfig, target_name: event.target.value })}
                  placeholder="Target name"
                />
                <button type="submit" disabled={loading}>
                  Export draft
                </button>
              </form>
              {exportResult && (
                <div className="exportResult">
                  <h3>Exported Files</h3>
                  <ul>
                    {exportResult.artifacts.map((artifact) => (
                      <li key={artifact.path}>
                        <strong>{artifact.description}:</strong> {artifact.path}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ) : (
            <p className="empty">Draft output will appear here after retrieval.</p>
          )}
        </section>

        <section className="panel wide">
          <h2>Indexed Sources</h2>
          <div className="sourceList">
            {sources.map((source) => (
              <article key={source.id} className="sourceCard">
                <div>
                  <strong>{source.title}</strong>
                  <p>{source.location}</p>
                </div>
                <span className="pill">{source.kind}</span>
              </article>
            ))}
            {sources.length === 0 && <p className="empty">Import a document or index repository specs to populate the corpus.</p>}
          </div>
        </section>
      </main>
    </div>
  );
}

async function assertOk(response) {
  const payload = await response.json();
  if (!response.ok) {
    throw new Error(payload.error ?? "Request failed");
  }
  return payload;
}

function slugify(input) {
  return String(input)
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}
