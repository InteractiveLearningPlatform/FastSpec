package app

import (
	"context"
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
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{})

	bundle, err := service.Search(context.Background(), "retrieval citations", 5)
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
	service := NewService(store, stubDOCXImporter{}, stubConfluenceImporter{}, stubIndexer{})

	draft, err := service.DraftSpec(context.Background(), "confluence drafting", "Speclist Draft", "openspec-markdown", 4)
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
