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
