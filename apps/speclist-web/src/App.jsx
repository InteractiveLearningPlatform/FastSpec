import { useEffect, useState } from "react";

const apiBase = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";

const initialConfluence = {
  base_url: "",
  page_id: "",
  token: "",
};

const sourceKinds = ["docx", "confluence", "spec"];
const draftPresets = ["general", "proposal", "design", "requirements"];
const reviewStatuses = ["ready", "needs-work", "blocked"];

const initialFilters = {
  kinds: [],
  origin: "",
  location_contains: "",
};

const initialExport = {
  mode: "filesystem",
  format: "openspec-markdown",
  target_dir: "./exports",
  target_name: "speclist-draft",
  change_name: "",
  artifact: "proposal",
  capability_name: "",
};

export default function App() {
  const [sources, setSources] = useState([]);
  const [changes, setChanges] = useState([]);
  const [searchQuery, setSearchQuery] = useState("spec retrieval citations");
  const [searchResults, setSearchResults] = useState([]);
  const [filters, setFilters] = useState(initialFilters);
  const [draftTitle, setDraftTitle] = useState("Speclist Draft");
  const [draftPreset, setDraftPreset] = useState("general");
  const [draft, setDraft] = useState(null);
  const [originalDraft, setOriginalDraft] = useState(null);
  const [reviewFlags, setReviewFlags] = useState({});
  const [citationInspection, setCitationInspection] = useState(null);
  const [sourceDetail, setSourceDetail] = useState(null);
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
      const [sourcesResponse, changesResponse] = await Promise.all([
        fetch(`${apiBase}/api/v1/sources`),
        fetch(`${apiBase}/api/v1/openspec/changes`),
      ]);
      const sourcesPayload = await sourcesResponse.json();
      const changesPayload = await changesResponse.json();
      setSources(sourcesPayload.sources ?? []);
      setChanges(changesPayload.changes ?? []);
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
        body: JSON.stringify({ query: searchQuery, limit: 8, filters }),
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
          filters,
          preset: draftPreset,
        }),
      });
      const payload = await assertOk(response);
      setDraft(payload);
      setOriginalDraft(structuredClone(payload));
      setReviewFlags({});
      setExportResult(null);
      setExportConfig((current) => ({
        ...current,
        target_name: slugify(payload.title || draftTitle),
        change_name: current.change_name || changes[0]?.name || "",
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
        body: JSON.stringify(buildExportPayload(draft, exportConfig)),
      });
      const payload = await assertOk(response);
      setExportResult(payload);
      setMessage(`Exported ${payload.artifacts?.length ?? 0} artifact(s)`);
    });
  }

  async function inspectCitation(citation) {
    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/citations/inspect`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ citation }),
      });
      const payload = await assertOk(response);
      setCitationInspection(payload);
      setMessage(`Loaded citation ${citation}`);
    });
  }

  async function inspectSource(sourceID) {
    await runAction(async () => {
      const response = await fetch(`${apiBase}/api/v1/sources/inspect`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ source_id: sourceID }),
      });
      const payload = await assertOk(response);
      setSourceDetail(payload.source);
      setMessage(`Loaded source ${payload.source.title}`);
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

  function updateDraft(updater) {
    setDraft((current) => (current ? updater(current) : current));
  }

  function updateReviewFlag(index, patch) {
    setReviewFlags((current) => ({
      ...current,
      [index]: { status: "ready", note: "", ...current[index], ...patch },
    }));
  }

  function removeSection(index) {
    updateDraft((current) => ({
      ...current,
      sections: current.sections.filter((_, currentIndex) => currentIndex !== index),
    }));
    setReviewFlags((current) => shiftFlagsAfterRemoval(current, index));
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
            <fieldset className="stack">
              <legend>Source filters</legend>
              <div className="sourceList">
                {sourceKinds.map((kind) => (
                  <label key={kind} className="sourceCard">
                    <span>{kind}</span>
                    <input
                      type="checkbox"
                      checked={filters.kinds.includes(kind)}
                      onChange={() =>
                        setFilters((current) => ({ ...current, kinds: toggleFilterKind(current.kinds, kind) }))
                      }
                    />
                  </label>
                ))}
              </div>
              <select value={filters.origin} onChange={(event) => setFilters({ ...filters, origin: event.target.value })}>
                <option value="">All origins</option>
                <option value="imported">Imported docs</option>
                <option value="repository">Repository specs</option>
              </select>
              <input
                value={filters.location_contains}
                onChange={(event) => setFilters({ ...filters, location_contains: event.target.value })}
                placeholder="Location contains"
              />
            </fieldset>
            <button type="submit" disabled={loading}>
              Search corpus
            </button>
          </form>
          <p className="empty">Applied filters: {describeFilters(filters)}</p>
          <div className="results">
            {searchResults.map((result) => (
              <article key={result.chunk.id} className="resultCard">
                <div className="resultMeta">
                  <strong>{result.source.title}</strong>
                  <span>{result.chunk.citation}</span>
                </div>
                <p>{result.chunk.text}</p>
                <button type="button" className="secondary" onClick={() => inspectCitation(result.chunk.citation)}>
                  Inspect citation
                </button>
              </article>
            ))}
            {searchResults.length === 0 && <p className="empty">No retrieval bundle loaded yet.</p>}
          </div>
        </section>

        <section className="panel">
          <h2>Create Draft</h2>
          <form className="stack" onSubmit={handleDraft}>
            <input value={draftTitle} onChange={(event) => setDraftTitle(event.target.value)} placeholder="Draft title" />
            <select value={draftPreset} onChange={(event) => setDraftPreset(event.target.value)}>
              {draftPresets.map((preset) => (
                <option key={preset} value={preset}>
                  {titleCase(preset)} preset
                </option>
              ))}
            </select>
            <button type="submit" disabled={loading}>
              Draft spec
            </button>
          </form>
          {draft ? (
            <div className="draft">
              <div className="stack">
                <p className="empty">Preset: {titleCase(draft.preset || "general")}</p>
                <input
                  value={draft.title}
                  onChange={(event) => updateDraft((current) => ({ ...current, title: event.target.value }))}
                  placeholder="Draft title"
                />
                <textarea
                  value={draft.summary}
                  onChange={(event) => updateDraft((current) => ({ ...current, summary: event.target.value }))}
                  rows={3}
                  placeholder="Draft summary"
                />
              </div>
              {draft.sections.map((section, index) => (
                <section key={`${section.heading}-${index}`} className="draftSection">
                  <div className="resultMeta">
                    <strong>Section {index + 1}</strong>
                    <button
                      type="button"
                      className="secondary"
                      onClick={() => removeSection(index)}
                    >
                      Remove
                    </button>
                  </div>
                  <div className="stack">
                    <select
                      value={reviewFlags[index]?.status || "ready"}
                      onChange={(event) => updateReviewFlag(index, { status: event.target.value })}
                    >
                      {reviewStatuses.map((status) => (
                        <option key={status} value={status}>
                          {titleCase(status)}
                        </option>
                      ))}
                    </select>
                    <textarea
                      value={reviewFlags[index]?.note || ""}
                      onChange={(event) => updateReviewFlag(index, { note: event.target.value })}
                      rows={2}
                      placeholder="Optional review note"
                    />
                  </div>
                  <input
                    value={section.heading}
                    onChange={(event) =>
                      updateDraft((current) => ({
                        ...current,
                        sections: updateSection(current.sections, index, { heading: event.target.value }),
                      }))
                    }
                    placeholder="Section heading"
                  />
                  <textarea
                    value={section.body}
                    onChange={(event) =>
                      updateDraft((current) => ({
                        ...current,
                        sections: updateSection(current.sections, index, { body: event.target.value }),
                      }))
                    }
                    rows={6}
                    placeholder="Section body"
                  />
                  <textarea
                    value={section.citations.join("\n")}
                    onChange={(event) =>
                      updateDraft((current) => ({
                        ...current,
                        sections: updateSection(current.sections, index, {
                          citations: splitCitations(event.target.value),
                        }),
                      }))
                    }
                    rows={3}
                    placeholder="One citation per line"
                  />
                  {section.citations.length > 0 && (
                    <div className="sourceList">
                      {section.citations.map((citation) => (
                        <button
                          key={`${citation}-${index}`}
                          type="button"
                          className="secondary"
                          onClick={() => inspectCitation(citation)}
                        >
                          Inspect {citation}
                        </button>
                      ))}
                    </div>
                  )}
                </section>
              ))}
              <button
                type="button"
                className="secondary"
                onClick={() =>
                  updateDraft((current) => ({
                    ...current,
                    sections: [...current.sections, { heading: "New Section", body: "", citations: [] }],
                  }))
                }
              >
                Add section
              </button>
              <div className="panel">
                <h3>Draft Diff</h3>
                {renderDraftDiff(originalDraft, draft)}
              </div>
              <div className="panel">
                <h3>Review Flags</h3>
                {renderReviewFlags(draft.sections, reviewFlags)}
              </div>
              <form className="stack exportForm" onSubmit={handleExport}>
                <h3>Export Reviewed Draft</h3>
                {renderExportReadiness(draft, originalDraft, reviewFlags)}
                <select
                  value={exportConfig.mode}
                  onChange={(event) => setExportConfig({ ...exportConfig, mode: event.target.value })}
                >
                  <option value="filesystem">Filesystem target</option>
                  <option value="openspec-change">OpenSpec change target</option>
                </select>
                <select
                  value={exportConfig.format}
                  onChange={(event) => setExportConfig({ ...exportConfig, format: event.target.value })}
                >
                  <option value="openspec-markdown">OpenSpec markdown</option>
                  <option value="fastspec-yaml">FastSpec YAML</option>
                </select>
                {exportConfig.mode === "filesystem" ? (
                  <>
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
                  </>
                ) : (
                  <>
                    <select
                      value={exportConfig.change_name}
                      onChange={(event) => setExportConfig({ ...exportConfig, change_name: event.target.value })}
                    >
                      <option value="">Select change</option>
                      {changes.map((change) => (
                        <option key={change.name} value={change.name}>
                          {change.name}
                        </option>
                      ))}
                    </select>
                    <select
                      value={exportConfig.artifact}
                      onChange={(event) => setExportConfig({ ...exportConfig, artifact: event.target.value })}
                    >
                      <option value="proposal">proposal.md</option>
                      <option value="design">design.md</option>
                      <option value="tasks">tasks.md</option>
                      <option value="spec">specs/&lt;capability&gt;/spec.md</option>
                    </select>
                    {exportConfig.artifact === "spec" && (
                      <input
                        value={exportConfig.capability_name}
                        onChange={(event) => setExportConfig({ ...exportConfig, capability_name: event.target.value })}
                        placeholder="Capability name"
                      />
                    )}
                  </>
                )}
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
                <div className="stack">
                  <span className="pill">{source.kind}</span>
                  <button type="button" className="secondary" onClick={() => inspectSource(source.id)}>
                    Inspect source
                  </button>
                </div>
              </article>
            ))}
            {sources.length === 0 && <p className="empty">Import a document or index repository specs to populate the corpus.</p>}
          </div>
        </section>

        <section className="panel wide">
          <div className="resultMeta">
            <h2>Citation Inspector</h2>
            {citationInspection && (
              <button type="button" className="secondary" onClick={() => setCitationInspection(null)}>
                Clear
              </button>
            )}
          </div>
          {citationInspection ? (
            <article className="resultCard">
              <div className="resultMeta">
                <strong>{citationInspection.source.title}</strong>
                <span>{citationInspection.citation}</span>
              </div>
              <p>{citationInspection.source.location}</p>
              <p>
                <strong>Section:</strong> {citationInspection.chunk.section}
              </p>
              <p>{citationInspection.chunk.text}</p>
              {Object.keys(citationInspection.source.metadata ?? {}).length > 0 && (
                <pre>{JSON.stringify(citationInspection.source.metadata, null, 2)}</pre>
              )}
            </article>
          ) : (
            <p className="empty">Inspect a citation from a search result or draft section to load its grounded source context.</p>
          )}
        </section>

        <section className="panel wide">
          <div className="resultMeta">
            <h2>Source Detail</h2>
            {sourceDetail && (
              <button type="button" className="secondary" onClick={() => setSourceDetail(null)}>
                Clear
              </button>
            )}
          </div>
          {sourceDetail ? (
            <article className="draft">
              <div className="resultMeta">
                <strong>{sourceDetail.title}</strong>
                <span>{sourceDetail.kind}</span>
              </div>
              <p>{sourceDetail.location}</p>
              {Object.keys(sourceDetail.metadata ?? {}).length > 0 && (
                <pre>{JSON.stringify(sourceDetail.metadata, null, 2)}</pre>
              )}
              <div className="results">
                {sourceDetail.chunks.map((chunk) => (
                  <article key={chunk.id} className="resultCard">
                    <div className="resultMeta">
                      <strong>{chunk.section}</strong>
                      <span>{chunk.citation}</span>
                    </div>
                    <p>{chunk.text}</p>
                    <button type="button" className="secondary" onClick={() => inspectCitation(chunk.citation)}>
                      Inspect citation
                    </button>
                  </article>
                ))}
              </div>
            </article>
          ) : (
            <p className="empty">Inspect a source from the indexed source list to review its metadata and chunk inventory.</p>
          )}
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

function buildExportPayload(draft, exportConfig) {
  if (exportConfig.mode === "openspec-change") {
    return {
      draft,
      format: "openspec-markdown",
      openspec_target: {
        change_name: exportConfig.change_name,
        artifact: exportConfig.artifact,
        capability_name: exportConfig.capability_name || undefined,
      },
    };
  }

  return {
    draft,
    format: exportConfig.format,
    target_dir: exportConfig.target_dir,
    target_name: exportConfig.target_name,
  };
}

function toggleFilterKind(currentKinds, kind) {
  return currentKinds.includes(kind) ? currentKinds.filter((entry) => entry !== kind) : [...currentKinds, kind];
}

function describeFilters(filters) {
  const parts = [];
  if (filters.kinds.length > 0) {
    parts.push(`kinds=${filters.kinds.join(", ")}`);
  }
  if (filters.origin) {
    parts.push(`origin=${filters.origin}`);
  }
  if (filters.location_contains.trim()) {
    parts.push(`location~${filters.location_contains.trim()}`);
  }
  return parts.length > 0 ? parts.join(" | ") : "none";
}

function updateSection(sections, index, patch) {
  return sections.map((section, sectionIndex) => (sectionIndex === index ? { ...section, ...patch } : section));
}

function splitCitations(value) {
  return value
    .split("\n")
    .map((citation) => citation.trim())
    .filter(Boolean);
}

function titleCase(input) {
  return String(input)
    .split(/[-_\s]+/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function renderDraftDiff(originalDraft, currentDraft) {
  if (!originalDraft || !currentDraft) {
    return <p className="empty">Generate a draft to inspect differences.</p>;
  }

  const sections = buildSectionDiffs(originalDraft.sections ?? [], currentDraft.sections ?? []);
  const hasTopLevelChanges =
    originalDraft.title !== currentDraft.title ||
    originalDraft.summary !== currentDraft.summary ||
    (originalDraft.preset || "general") !== (currentDraft.preset || "general");
  const hasSectionChanges = sections.some((section) => section.kind !== "unchanged");

  if (!hasTopLevelChanges && !hasSectionChanges) {
    return <p className="empty">No edits yet. The current draft still matches the generated draft.</p>;
  }

  return (
    <div className="stack">
      {originalDraft.title !== currentDraft.title && (
        <article className="resultCard">
          <strong>Title changed</strong>
          <p>Original: {originalDraft.title}</p>
          <p>Current: {currentDraft.title}</p>
        </article>
      )}
      {originalDraft.summary !== currentDraft.summary && (
        <article className="resultCard">
          <strong>Summary changed</strong>
          <p>Original: {originalDraft.summary}</p>
          <p>Current: {currentDraft.summary}</p>
        </article>
      )}
      {(originalDraft.preset || "general") !== (currentDraft.preset || "general") && (
        <article className="resultCard">
          <strong>Preset changed</strong>
          <p>Original: {titleCase(originalDraft.preset || "general")}</p>
          <p>Current: {titleCase(currentDraft.preset || "general")}</p>
        </article>
      )}
      {sections
        .filter((section) => section.kind !== "unchanged")
        .map((section) => (
          <article key={section.key} className="resultCard">
            <strong>{section.label}</strong>
            {section.kind === "added" ? (
              <>
                <p>Heading: {section.current.heading}</p>
                <p>Body: {section.current.body}</p>
              </>
            ) : section.kind === "removed" ? (
              <>
                <p>Heading: {section.original.heading}</p>
                <p>Body: {section.original.body}</p>
              </>
            ) : (
              <>
                {section.original.heading !== section.current.heading && (
                  <p>
                    Heading: {section.original.heading} {"->"} {section.current.heading}
                  </p>
                )}
                {section.original.body !== section.current.body && (
                  <p>
                    Body changed from:
                    {" "}
                    {clipForDiff(section.original.body)}
                    {" "}
                    {"->"} {clipForDiff(section.current.body)}
                  </p>
                )}
                {joinCitations(section.original.citations) !== joinCitations(section.current.citations) && (
                  <p>
                    Citations: {joinCitations(section.original.citations)} {"->"} {joinCitations(section.current.citations)}
                  </p>
                )}
              </>
            )}
          </article>
        ))}
    </div>
  );
}

function buildSectionDiffs(originalSections, currentSections) {
  const length = Math.max(originalSections.length, currentSections.length);
  const diffs = [];
  for (let index = 0; index < length; index += 1) {
    const original = originalSections[index];
    const current = currentSections[index];
    if (!original && current) {
      diffs.push({ key: `section-${index}`, label: `Section ${index + 1} added`, kind: "added", current });
      continue;
    }
    if (original && !current) {
      diffs.push({ key: `section-${index}`, label: `Section ${index + 1} removed`, kind: "removed", original });
      continue;
    }
    const changed =
      original.heading !== current.heading ||
      original.body !== current.body ||
      joinCitations(original.citations) !== joinCitations(current.citations);
    diffs.push({
      key: `section-${index}`,
      label: changed ? `Section ${index + 1} changed` : `Section ${index + 1}`,
      kind: changed ? "changed" : "unchanged",
      original,
      current,
    });
  }
  return diffs;
}

function joinCitations(citations = []) {
  return citations.join(" | ");
}

function clipForDiff(input) {
  const value = String(input || "").trim();
  return value.length > 120 ? `${value.slice(0, 117)}...` : value || "(empty)";
}

function renderReviewFlags(sections, reviewFlags) {
  const flaggedSections = sections
    .map((section, index) => ({ index, section, flag: reviewFlags[index] || { status: "ready", note: "" } }))
    .filter(({ flag }) => flag.status !== "ready" || flag.note.trim());

  if (flaggedSections.length === 0) {
    return <p className="empty">No follow-up flags. All sections are currently marked ready.</p>;
  }

  return (
    <div className="stack">
      {flaggedSections.map(({ index, section, flag }) => (
        <article key={`flag-${index}`} className="resultCard">
          <strong>
            Section {index + 1}: {section.heading}
          </strong>
          <p>Status: {titleCase(flag.status)}</p>
          {flag.note.trim() && <p>Note: {flag.note.trim()}</p>}
        </article>
      ))}
    </div>
  );
}

function shiftFlagsAfterRemoval(flags, removedIndex) {
  const shifted = {};
  Object.entries(flags).forEach(([key, value]) => {
    const index = Number(key);
    if (index < removedIndex) {
      shifted[index] = value;
    } else if (index > removedIndex) {
      shifted[index - 1] = value;
    }
  });
  return shifted;
}

function renderExportReadiness(draft, originalDraft, reviewFlags) {
  const readiness = computeExportReadiness(draft, originalDraft, reviewFlags);
  return (
    <div className="resultCard">
      <strong>Export Readiness: {titleCase(readiness.status)}</strong>
      {readiness.blockers.length > 0 && (
        <div>
          <p>Blockers</p>
          <ul>
            {readiness.blockers.map((item) => (
              <li key={`blocker-${item}`}>{item}</li>
            ))}
          </ul>
        </div>
      )}
      {readiness.warnings.length > 0 && (
        <div>
          <p>Warnings</p>
          <ul>
            {readiness.warnings.map((item) => (
              <li key={`warning-${item}`}>{item}</li>
            ))}
          </ul>
        </div>
      )}
      {readiness.blockers.length === 0 && readiness.warnings.length === 0 && (
        <p className="empty">No blockers or warnings. The draft is ready for export.</p>
      )}
    </div>
  );
}

function computeExportReadiness(draft, originalDraft, reviewFlags) {
  const blockers = [];
  const warnings = [];

  if (!String(draft?.title || "").trim()) {
    blockers.push("Draft title is empty.");
  }
  if (!String(draft?.summary || "").trim()) {
    blockers.push("Draft summary is empty.");
  }
  if (!draft?.sections?.length) {
    blockers.push("Draft has no sections.");
  }

  (draft?.sections || []).forEach((section, index) => {
    if (!String(section.heading || "").trim()) {
      blockers.push(`Section ${index + 1} is missing a heading.`);
    }
    if (!String(section.body || "").trim()) {
      blockers.push(`Section ${index + 1} is missing body content.`);
    }
    if ((section.citations || []).length === 0) {
      warnings.push(`Section ${index + 1} has no citations.`);
    }

    const flag = reviewFlags[index];
    if (!flag) {
      return;
    }
    if (flag.status === "blocked") {
      blockers.push(`Section ${index + 1} is marked blocked.`);
    }
    if (flag.status === "needs-work") {
      warnings.push(`Section ${index + 1} is marked needs-work.`);
    }
    if (String(flag.note || "").trim()) {
      warnings.push(`Section ${index + 1} has a review note.`);
    }
  });

  if (originalDraft && JSON.stringify(originalDraft) === JSON.stringify(draft)) {
    warnings.push("Draft is unchanged from the originally generated version.");
  }

  const status = blockers.length > 0 ? "blocked" : warnings.length > 0 ? "warning" : "ready";
  return { status, blockers: uniqueItems(blockers), warnings: uniqueItems(warnings) };
}

function uniqueItems(items) {
  return [...new Set(items)];
}
