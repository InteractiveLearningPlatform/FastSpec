package domain

import (
	"context"
	"time"
)

type SourceKind string

const (
	SourceKindDocx       SourceKind = "docx"
	SourceKindConfluence SourceKind = "confluence"
	SourceKindSpec       SourceKind = "spec"
)

type Chunk struct {
	ID       string            `json:"id"`
	SourceID string            `json:"source_id"`
	Section  string            `json:"section"`
	Text     string            `json:"text"`
	Citation string            `json:"citation"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type SourceDocument struct {
	ID         string            `json:"id"`
	Kind       SourceKind        `json:"kind"`
	Title      string            `json:"title"`
	Location   string            `json:"location"`
	ImportedAt time.Time         `json:"imported_at"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Chunks     []Chunk           `json:"chunks"`
}

type RetrievalResult struct {
	Chunk  Chunk      `json:"chunk"`
	Source SourceStub `json:"source"`
	Score  int        `json:"score"`
}

type SourceStub struct {
	ID       string            `json:"id"`
	Kind     SourceKind        `json:"kind"`
	Title    string            `json:"title"`
	Location string            `json:"location"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type RetrievalBundle struct {
	Query   string            `json:"query"`
	Results []RetrievalResult `json:"results"`
}

type DraftSection struct {
	Heading   string   `json:"heading"`
	Body      string   `json:"body"`
	Citations []string `json:"citations"`
}

type DraftSpec struct {
	Title       string         `json:"title"`
	Query       string         `json:"query"`
	Format      string         `json:"format"`
	Summary     string         `json:"summary"`
	Sections    []DraftSection `json:"sections"`
	SourceCount int            `json:"source_count"`
}

type CorpusStore interface {
	Save(ctx context.Context, document SourceDocument) error
	List(ctx context.Context) ([]SourceDocument, error)
}

type DocxImporter interface {
	Import(ctx context.Context, filename string, contents []byte) (SourceDocument, error)
}

type ConfluenceImporter interface {
	Import(ctx context.Context, request ConfluenceImportRequest) (SourceDocument, error)
}

type SpecIndexer interface {
	Index(ctx context.Context, repoRoot string) ([]SourceDocument, error)
}

type ConfluenceImportRequest struct {
	BaseURL string `json:"base_url"`
	PageID  string `json:"page_id"`
	Token   string `json:"token,omitempty"`
}
