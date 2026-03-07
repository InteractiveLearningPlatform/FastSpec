package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/app"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type emptyStore struct{}

func (emptyStore) Save(_ context.Context, _ domain.SourceDocument) error   { return nil }
func (emptyStore) List(_ context.Context) ([]domain.SourceDocument, error) { return nil, nil }

type emptyDOCX struct{}

func (emptyDOCX) Import(_ context.Context, _ string, _ []byte) (domain.SourceDocument, error) {
	return domain.SourceDocument{}, nil
}

type emptyConfluence struct{}

func (emptyConfluence) Import(_ context.Context, _ domain.ConfluenceImportRequest) (domain.SourceDocument, error) {
	return domain.SourceDocument{}, nil
}

type emptyIndexer struct{}

func (emptyIndexer) Index(_ context.Context, _ string) ([]domain.SourceDocument, error) {
	return nil, nil
}

type seededStore struct {
	documents []domain.SourceDocument
}

func (s seededStore) Save(_ context.Context, _ domain.SourceDocument) error {
	return nil
}

func (s seededStore) List(_ context.Context) ([]domain.SourceDocument, error) {
	return s.documents, nil
}

func TestExportEndpointWritesArtifacts(t *testing.T) {
	service := app.NewService(emptyStore{}, emptyDOCX{}, emptyConfluence{}, emptyIndexer{}, "")
	server := NewServer(service, "")

	body, err := json.Marshal(domain.DraftExportRequest{
		Draft: domain.DraftSpec{
			Title:       "Speclist Export",
			Query:       "export reviewed drafts",
			Summary:     "Grounded export draft",
			SourceCount: 1,
			Sections: []domain.DraftSection{
				{Heading: "Why", Body: "Need durable output.", Citations: []string{"notes.docx > Why"}},
			},
		},
		Format:     domain.ExportFormatOpenSpecMarkdown,
		TargetDir:  t.TempDir(),
		TargetName: "exported-draft",
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/exports", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("unexpected status %d: %s", recorder.Code, recorder.Body.String())
	}

	var result domain.DraftExportResult
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result.Artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(result.Artifacts))
	}
	if filepath.Ext(result.Artifacts[0].Path) != ".md" {
		t.Fatalf("expected markdown export, got %s", result.Artifacts[0].Path)
	}
}

func TestSearchEndpointAppliesFilters(t *testing.T) {
	service := app.NewService(seededStore{
		documents: []domain.SourceDocument{
			{
				ID:       "spec-1",
				Kind:     domain.SourceKindSpec,
				Title:    "Repository Search",
				Location: "openspec/specs/search.md",
				Chunks: []domain.Chunk{
					{ID: "chunk-1", SourceID: "spec-1", Section: "Requirement", Text: "Repository search should stay searchable.", Citation: "search.md > Requirement"},
				},
			},
			{
				ID:       "doc-1",
				Kind:     domain.SourceKindDocx,
				Title:    "Imported Search Notes",
				Location: "notes/search.docx",
				Chunks: []domain.Chunk{
					{ID: "chunk-2", SourceID: "doc-1", Section: "Notes", Text: "Imported notes mention search.", Citation: "notes.docx > Notes"},
				},
			},
		},
	}, emptyDOCX{}, emptyConfluence{}, emptyIndexer{}, "")
	server := NewServer(service, "")

	body, err := json.Marshal(map[string]any{
		"query": "search searchable",
		"limit": 5,
		"filters": map[string]any{
			"origin":            "repository",
			"location_contains": "openspec/specs",
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/search", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", recorder.Code, recorder.Body.String())
	}

	var result domain.RetrievalBundle
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 filtered result, got %d", len(result.Results))
	}
	if result.Results[0].Source.Kind != domain.SourceKindSpec {
		t.Fatalf("expected spec result, got %s", result.Results[0].Source.Kind)
	}
}

func TestInspectCitationEndpointReturnsSourceContext(t *testing.T) {
	service := app.NewService(seededStore{
		documents: []domain.SourceDocument{
			{
				ID:       "doc-1",
				Kind:     domain.SourceKindDocx,
				Title:    "Imported Search Notes",
				Location: "notes/search.docx",
				Chunks: []domain.Chunk{
					{ID: "chunk-2", SourceID: "doc-1", Section: "Notes", Text: "Imported notes mention search.", Citation: "notes.docx > Notes"},
				},
			},
		},
	}, emptyDOCX{}, emptyConfluence{}, emptyIndexer{}, "")
	server := NewServer(service, "")

	body, err := json.Marshal(map[string]any{"citation": "notes.docx > Notes"})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/citations/inspect", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", recorder.Code, recorder.Body.String())
	}

	var inspection domain.CitationInspection
	if err := json.Unmarshal(recorder.Body.Bytes(), &inspection); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if inspection.Chunk.Section != "Notes" {
		t.Fatalf("expected notes section, got %s", inspection.Chunk.Section)
	}
	if inspection.Source.Location != "notes/search.docx" {
		t.Fatalf("expected source location, got %s", inspection.Source.Location)
	}
}
