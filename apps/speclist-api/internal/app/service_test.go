package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type memoryStore struct {
	documents []domain.SourceDocument
}

func (m *memoryStore) Save(_ context.Context, document domain.SourceDocument) error {
	m.documents = append(m.documents, document)
	return nil
}

func (m *memoryStore) List(_ context.Context) ([]domain.SourceDocument, error) {
	return m.documents, nil
}

type stubDOCXImporter struct {
	document domain.SourceDocument
}

func (s stubDOCXImporter) Import(_ context.Context, _ string, _ []byte) (domain.SourceDocument, error) {
	return s.document, nil
}

type stubConfluenceImporter struct {
	document domain.SourceDocument
}

func (s stubConfluenceImporter) Import(_ context.Context, _ domain.ConfluenceImportRequest) (domain.SourceDocument, error) {
	return s.document, nil
}

type stubIndexer struct {
	documents []domain.SourceDocument
}

func (s stubIndexer) Index(_ context.Context, _ string) ([]domain.SourceDocument, error) {
	return s.documents, nil
}

func TestSearchReturnsGroundedResults(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "spec-1",
				Kind:       domain.SourceKindSpec,
				Title:      "RAG Search",
				Location:   "openspec/specs/rag.md",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "spec-1", Section: "Requirement", Text: "System MUST support compact retrieval bundles with citations.", Citation: "rag.md > Requirement"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	bundle, err := service.Search(context.Background(), "retrieval citations", 5, domain.RetrievalFilter{})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(bundle.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(bundle.Results))
	}
	if bundle.Results[0].Chunk.Citation != "rag.md > Requirement" {
		t.Fatalf("unexpected citation: %s", bundle.Results[0].Chunk.Citation)
	}
}

func TestDraftSpecIncludesCitations(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "doc-1",
				Kind:       domain.SourceKindDocx,
				Title:      "Platform Notes",
				Location:   "notes.docx",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "doc-1", Section: "Goals", Text: "Support confluence ingestion and spec drafting.", Citation: "notes.docx > Goals"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	draft, err := service.DraftSpec(context.Background(), "confluence drafting", "Speclist Draft", "openspec-markdown", 4, domain.RetrievalFilter{}, domain.DraftPresetGeneral)
	if err != nil {
		t.Fatalf("draft failed: %v", err)
	}
	if len(draft.Sections) == 0 {
		t.Fatal("expected draft sections")
	}
	if len(draft.Sections[0].Citations) == 0 {
		t.Fatal("expected citations in draft")
	}
}

func TestSearchAppliesSourceFilters(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "spec-1",
				Kind:       domain.SourceKindSpec,
				Title:      "Repository Retrieval",
				Location:   "openspec/specs/retrieval.md",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "spec-1", Section: "Requirement", Text: "Repository specs SHOULD stay searchable.", Citation: "retrieval.md > Requirement"},
				},
			},
			{
				ID:         "doc-1",
				Kind:       domain.SourceKindDocx,
				Title:      "Imported Notes",
				Location:   "notes/retrieval-guidance.docx",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-2", SourceID: "doc-1", Section: "Notes", Text: "Imported notes also mention retrieval.", Citation: "notes.docx > Notes"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	bundle, err := service.Search(context.Background(), "retrieval searchable", 5, domain.RetrievalFilter{
		Origin:           domain.SourceOriginRepository,
		LocationContains: "openspec/specs",
	})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(bundle.Results) != 1 {
		t.Fatalf("expected 1 filtered result, got %d", len(bundle.Results))
	}
	if bundle.Results[0].Source.Kind != domain.SourceKindSpec {
		t.Fatalf("expected spec source, got %s", bundle.Results[0].Source.Kind)
	}
	if bundle.Filters.Origin != domain.SourceOriginRepository {
		t.Fatalf("expected repository origin, got %s", bundle.Filters.Origin)
	}
}

func TestDraftSpecReusesSearchFilters(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "docx-1",
				Kind:       domain.SourceKindDocx,
				Title:      "Imported Product Notes",
				Location:   "product/notes.docx",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "docx-1", Section: "Goals", Text: "Imported notes define the product scope.", Citation: "notes.docx > Goals"},
				},
			},
			{
				ID:         "spec-1",
				Kind:       domain.SourceKindSpec,
				Title:      "Repository Plan",
				Location:   "openspec/specs/plan.md",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-2", SourceID: "spec-1", Section: "Context", Text: "Repository specs should not appear here.", Citation: "plan.md > Context"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	draft, err := service.DraftSpec(context.Background(), "product scope", "Imported Draft", "openspec-markdown", 4, domain.RetrievalFilter{
		Kinds:  []domain.SourceKind{domain.SourceKindDocx},
		Origin: domain.SourceOriginImported,
	}, domain.DraftPresetGeneral)
	if err != nil {
		t.Fatalf("draft failed: %v", err)
	}
	if draft.SourceCount != 1 {
		t.Fatalf("expected 1 filtered source result, got %d", draft.SourceCount)
	}
	if draft.Filters.Origin != domain.SourceOriginImported {
		t.Fatalf("expected imported origin, got %s", draft.Filters.Origin)
	}
	if len(draft.Filters.Kinds) != 1 || draft.Filters.Kinds[0] != domain.SourceKindDocx {
		t.Fatalf("expected docx filter, got %+v", draft.Filters.Kinds)
	}
}

func TestDraftSpecSupportsProposalPreset(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "doc-1",
				Kind:       domain.SourceKindDocx,
				Title:      "Platform Notes",
				Location:   "notes.docx",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "doc-1", Section: "Goals", Text: "Support proposal-oriented drafting.", Citation: "notes.docx > Goals"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	draft, err := service.DraftSpec(context.Background(), "proposal drafting", "Proposal Draft", "openspec-markdown", 4, domain.RetrievalFilter{}, domain.DraftPresetProposal)
	if err != nil {
		t.Fatalf("draft failed: %v", err)
	}
	if draft.Preset != domain.DraftPresetProposal {
		t.Fatalf("expected proposal preset, got %s", draft.Preset)
	}
	if len(draft.Sections) != 3 || draft.Sections[1].Heading != "What Changes" {
		t.Fatalf("unexpected proposal sections: %+v", draft.Sections)
	}
}

func TestDraftSpecSupportsRequirementsPreset(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "doc-1",
				Kind:       domain.SourceKindDocx,
				Title:      "Platform Notes",
				Location:   "notes.docx",
				ImportedAt: time.Now().UTC(),
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "doc-1", Section: "Goals", Text: "Support requirement-oriented drafting.", Citation: "notes.docx > Goals"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	draft, err := service.DraftSpec(context.Background(), "requirements drafting", "Requirements Draft", "openspec-markdown", 4, domain.RetrievalFilter{}, domain.DraftPresetRequirements)
	if err != nil {
		t.Fatalf("draft failed: %v", err)
	}
	if len(draft.Sections) != 3 || draft.Sections[1].Heading != "Requirements" || draft.Sections[2].Heading != "Scenarios" {
		t.Fatalf("unexpected requirements sections: %+v", draft.Sections)
	}
}

func TestInspectCitationReturnsSourceContext(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "doc-1",
				Kind:       domain.SourceKindDocx,
				Title:      "Imported Notes",
				Location:   "notes/product.docx",
				ImportedAt: time.Now().UTC(),
				Metadata:   map[string]string{"team": "platform"},
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "doc-1", Section: "Goals", Text: "Support grounded draft review.", Citation: "product.docx > Goals"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	inspection, err := service.InspectCitation(context.Background(), "product.docx > Goals")
	if err != nil {
		t.Fatalf("inspect citation failed: %v", err)
	}
	if inspection.Source.Title != "Imported Notes" {
		t.Fatalf("unexpected source title: %s", inspection.Source.Title)
	}
	if inspection.Chunk.Text != "Support grounded draft review." {
		t.Fatalf("unexpected chunk text: %s", inspection.Chunk.Text)
	}
}

func TestInspectSourceReturnsFullDocument(t *testing.T) {
	store := &memoryStore{
		documents: []domain.SourceDocument{
			{
				ID:         "doc-1",
				Kind:       domain.SourceKindDocx,
				Title:      "Imported Notes",
				Location:   "notes/product.docx",
				ImportedAt: time.Now().UTC(),
				Metadata:   map[string]string{"team": "platform"},
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "doc-1", Section: "Goals", Text: "Support grounded draft review.", Citation: "product.docx > Goals"},
				},
			},
		},
	}
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	document, err := service.InspectSource(context.Background(), "doc-1")
	if err != nil {
		t.Fatalf("inspect source failed: %v", err)
	}
	if document.Title != "Imported Notes" {
		t.Fatalf("unexpected source title: %s", document.Title)
	}
	if len(document.Chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(document.Chunks))
	}
}

func TestInspectSourceRejectsMissingID(t *testing.T) {
	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	_, err := service.InspectSource(context.Background(), "missing-source")
	if err == nil || !strings.Contains(err.Error(), "was not found") {
		t.Fatalf("expected missing source error, got %v", err)
	}
}

func TestInspectCitationRejectsMissingCitation(t *testing.T) {
	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	_, err := service.InspectCitation(context.Background(), "missing > Citation")
	if err == nil || !strings.Contains(err.Error(), "was not found") {
		t.Fatalf("expected missing citation error, got %v", err)
	}
}

func TestExportDraftWritesMarkdownAndSidecar(t *testing.T) {
	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")
	targetDir := t.TempDir()

	result, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Exported Draft",
			Query:       "reviewed draft export",
			Summary:     "Draft summary",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "Why", Body: "Need a durable artifact.", Citations: []string{"notes.docx > Why"}},
			},
		},
		Format:     domain.ExportFormatOpenSpecMarkdown,
		TargetDir:  targetDir,
		TargetName: "exported-draft",
	})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if len(result.Artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(result.Artifacts))
	}
	contents, err := os.ReadFile(filepath.Join(targetDir, "exported-draft.md"))
	if err != nil {
		t.Fatalf("read exported markdown: %v", err)
	}
	if string(contents) == "" {
		t.Fatal("expected markdown content")
	}
}

func TestExportDraftWritesTypedFastSpecYAML(t *testing.T) {
	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")
	targetDir := t.TempDir()

	_, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Structured YAML Draft",
			Query:       "typed yaml export",
			Summary:     "Draft summary",
			SourceCount: 2,
			Sections: []domain.DraftSection{
				{Heading: "Why", Body: "Need stronger YAML export.", Citations: []string{"notes.docx > Why"}},
				{Heading: "Context", Body: "Support durable spec refinement.", Citations: []string{"notes.docx > Context"}},
				{Heading: "Proposed Requirements", Body: "- MUST reflect: Preserve requirement structure"},
			},
		},
		Format:     domain.ExportFormatFastSpecYAML,
		TargetDir:  targetDir,
		TargetName: "structured-yaml-draft",
	})
	if err != nil {
		t.Fatalf("yaml export failed: %v", err)
	}

	contents, err := os.ReadFile(filepath.Join(targetDir, "structured-yaml-draft.fastspec.yaml"))
	if err != nil {
		t.Fatalf("read exported yaml: %v", err)
	}
	text := string(contents)
	if !strings.Contains(text, "kind: SpecDocumentDraft") {
		t.Fatalf("expected typed draft kind, got: %s", text)
	}
	if !strings.Contains(text, "requirements:") || !strings.Contains(text, "statement:") {
		t.Fatalf("expected structured requirements, got: %s", text)
	}
}

func TestExportDraftRejectsOverwrite(t *testing.T) {
	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")
	targetDir := t.TempDir()
	path := filepath.Join(targetDir, "existing.md")
	if err := os.WriteFile(path, []byte("occupied"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	_, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Existing Draft",
			Query:       "overwrite check",
			Summary:     "Draft summary",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "Why", Body: "Need a durable artifact.", Citations: []string{"notes.docx > Why"}},
			},
		},
		Format:     domain.ExportFormatOpenSpecMarkdown,
		TargetDir:  targetDir,
		TargetName: "existing",
	})
	if err == nil {
		t.Fatal("expected overwrite rejection")
	}
}

func TestExportDraftRejectsBlankSectionHeading(t *testing.T) {
	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")

	_, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Edited Draft",
			Query:       "reviewed editing",
			Summary:     "Draft summary",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "   ", Body: "Need a durable artifact.", Citations: []string{"notes.docx > Why"}},
			},
		},
		Format:     domain.ExportFormatOpenSpecMarkdown,
		TargetDir:  t.TempDir(),
		TargetName: "edited-draft",
	})
	if err == nil || !strings.Contains(err.Error(), "section heading") {
		t.Fatalf("expected section heading validation error, got %v", err)
	}
}

func TestExportDraftNormalizesEditedContent(t *testing.T) {
	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, "")
	targetDir := t.TempDir()

	_, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "  Edited Draft  ",
			Query:       "  reviewed editing  ",
			Summary:     "  Draft summary  ",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "  Why  ", Body: "  Need a durable artifact.  ", Citations: []string{"  notes.docx > Why  ", "   "}},
			},
		},
		Format:     domain.ExportFormatOpenSpecMarkdown,
		TargetDir:  targetDir,
		TargetName: "edited-draft",
	})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	contents, err := os.ReadFile(filepath.Join(targetDir, "edited-draft.md"))
	if err != nil {
		t.Fatalf("read edited markdown: %v", err)
	}
	text := string(contents)
	if strings.Contains(text, "  Edited Draft  ") {
		t.Fatalf("expected normalized title, got: %s", text)
	}
	if !strings.Contains(text, "# Edited Draft") {
		t.Fatalf("expected trimmed title, got: %s", text)
	}
	if !strings.Contains(text, "- notes.docx > Why") {
		t.Fatalf("expected trimmed citation, got: %s", text)
	}
}

func TestExportDraftWritesToOpenSpecChangeTarget(t *testing.T) {
	repoRoot := t.TempDir()
	changeDir := filepath.Join(repoRoot, "openspec", "changes", "demo-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("mkdir change dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, ".openspec.yaml"), []byte("schema: spec-driven\n"), 0o644); err != nil {
		t.Fatalf("write change marker: %v", err)
	}

	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, repoRoot)
	result, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Proposal Draft",
			Query:       "change target export",
			Summary:     "Proposal summary",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "Why", Body: "Need to export into an active change.", Citations: []string{"notes.docx > Why"}},
			},
		},
		Format: domain.ExportFormatOpenSpecMarkdown,
		OpenSpecTarget: &domain.OpenSpecExportTarget{
			ChangeName: "demo-change",
			Artifact:   "proposal",
		},
	})
	if err != nil {
		t.Fatalf("export to change target failed: %v", err)
	}
	if len(result.Artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(result.Artifacts))
	}
	contents, err := os.ReadFile(filepath.Join(changeDir, "proposal.md"))
	if err != nil {
		t.Fatalf("expected proposal export: %v", err)
	}
	if !strings.Contains(string(contents), "## Why") || !strings.Contains(string(contents), "## What Changes") {
		t.Fatalf("expected proposal template, got: %s", string(contents))
	}
}

func TestListOpenSpecChangesReturnsActiveChanges(t *testing.T) {
	repoRoot := t.TempDir()
	active := filepath.Join(repoRoot, "openspec", "changes", "demo-change")
	archived := filepath.Join(repoRoot, "openspec", "changes", "archive", "old-change")
	if err := os.MkdirAll(active, 0o755); err != nil {
		t.Fatalf("mkdir active dir: %v", err)
	}
	if err := os.MkdirAll(archived, 0o755); err != nil {
		t.Fatalf("mkdir archived dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(active, ".openspec.yaml"), []byte("schema: spec-driven\n"), 0o644); err != nil {
		t.Fatalf("write active marker: %v", err)
	}
	if err := os.WriteFile(filepath.Join(archived, ".openspec.yaml"), []byte("schema: spec-driven\n"), 0o644); err != nil {
		t.Fatalf("write archive marker: %v", err)
	}

	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, repoRoot)
	changes, err := service.ListOpenSpecChanges(context.Background())
	if err != nil {
		t.Fatalf("list changes failed: %v", err)
	}
	if len(changes) != 1 || changes[0].Name != "demo-change" {
		t.Fatalf("unexpected changes: %+v", changes)
	}
}

func TestExportDraftRendersTasksTemplate(t *testing.T) {
	repoRoot := t.TempDir()
	changeDir := filepath.Join(repoRoot, "openspec", "changes", "demo-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("mkdir change dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, ".openspec.yaml"), []byte("schema: spec-driven\n"), 0o644); err != nil {
		t.Fatalf("write change marker: %v", err)
	}

	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, repoRoot)
	_, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Task Draft",
			Query:       "task rendering",
			Summary:     "Task summary",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "Proposed Requirements", Body: "- MUST reflect: Add indexing\n- MUST reflect: Add review flow"},
			},
		},
		Format: domain.ExportFormatOpenSpecMarkdown,
		OpenSpecTarget: &domain.OpenSpecExportTarget{
			ChangeName: "demo-change",
			Artifact:   "tasks",
		},
	})
	if err != nil {
		t.Fatalf("export tasks failed: %v", err)
	}
	contents, err := os.ReadFile(filepath.Join(changeDir, "tasks.md"))
	if err != nil {
		t.Fatalf("read tasks export: %v", err)
	}
	if !strings.Contains(string(contents), "- [ ] 1.1") {
		t.Fatalf("expected checklist template, got: %s", string(contents))
	}
}

func TestExportDraftRendersSpecTemplate(t *testing.T) {
	repoRoot := t.TempDir()
	changeDir := filepath.Join(repoRoot, "openspec", "changes", "demo-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("mkdir change dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, ".openspec.yaml"), []byte("schema: spec-driven\n"), 0o644); err != nil {
		t.Fatalf("write change marker: %v", err)
	}

	service := NewService(&memoryStore{}, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{}, repoRoot)
	_, err := service.ExportDraft(context.Background(), domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Spec Draft",
			Query:       "spec rendering",
			Summary:     "Spec summary",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "Why", Body: "Need requirement-oriented export."},
				{Heading: "Proposed Requirements", Body: "- MUST reflect: Keep structured scenarios"},
			},
		},
		Format: domain.ExportFormatOpenSpecMarkdown,
		OpenSpecTarget: &domain.OpenSpecExportTarget{
			ChangeName:     "demo-change",
			Artifact:       "spec",
			CapabilityName: "typed-rendering",
		},
	})
	if err != nil {
		t.Fatalf("export spec failed: %v", err)
	}
	contents, err := os.ReadFile(filepath.Join(changeDir, "specs", "typed-rendering", "spec.md"))
	if err != nil {
		t.Fatalf("read spec export: %v", err)
	}
	if !strings.Contains(string(contents), "### Requirement: Typed rendering") {
		t.Fatalf("expected spec template, got: %s", string(contents))
	}
}
