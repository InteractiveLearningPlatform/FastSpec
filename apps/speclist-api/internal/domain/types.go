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
	SourceKindMarkdown   SourceKind = "markdown"
	SourceKindCode       SourceKind = "code"
	SourceKindIR         SourceKind = "ir"
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

type SourceOrigin string

const (
	SourceOriginImported   SourceOrigin = "imported"
	SourceOriginRepository SourceOrigin = "repository"
)

type RetrievalFilter struct {
	Kinds            []SourceKind `json:"kinds,omitempty"`
	Origin           SourceOrigin `json:"origin,omitempty"`
	LocationContains string       `json:"location_contains,omitempty"`
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
	Filters RetrievalFilter   `json:"filters,omitempty"`
	Results []RetrievalResult `json:"results"`
}

type DraftSection struct {
	Heading   string   `json:"heading"`
	Body      string   `json:"body"`
	Citations []string `json:"citations"`
}

type DraftPreset string

const (
	DraftPresetGeneral      DraftPreset = "general"
	DraftPresetProposal     DraftPreset = "proposal"
	DraftPresetDesign       DraftPreset = "design"
	DraftPresetRequirements DraftPreset = "requirements"
)

type DraftSpec struct {
	Title       string          `json:"title"`
	Query       string          `json:"query"`
	Filters     RetrievalFilter `json:"filters,omitempty"`
	Preset      DraftPreset     `json:"preset,omitempty"`
	Format      string          `json:"format"`
	Summary     string          `json:"summary"`
	Sections    []DraftSection  `json:"sections"`
	SourceCount int             `json:"source_count"`
}

type CitationInspection struct {
	Citation string     `json:"citation"`
	Source   SourceStub `json:"source"`
	Chunk    Chunk      `json:"chunk"`
}

type ExportFormat string

const (
	ExportFormatOpenSpecMarkdown ExportFormat = "openspec-markdown"
	ExportFormatFastSpecYAML     ExportFormat = "fastspec-yaml"
)

type DraftExportRequest struct {
	Draft          DraftSpec             `json:"draft"`
	Format         ExportFormat          `json:"format"`
	TargetDir      string                `json:"target_dir"`
	TargetName     string                `json:"target_name"`
	OpenSpecTarget *OpenSpecExportTarget `json:"openspec_target,omitempty"`
}

type ExportArtifact struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

type DraftExportResult struct {
	Format    ExportFormat     `json:"format"`
	Artifacts []ExportArtifact `json:"artifacts"`
}

type OpenSpecExportTarget struct {
	ChangeName     string `json:"change_name"`
	Artifact       string `json:"artifact"`
	CapabilityName string `json:"capability_name,omitempty"`
}

type OpenSpecChange struct {
	Name      string   `json:"name"`
	Artifacts []string `json:"artifacts"`
}

type AssetKind string

const (
	AssetKindPublishedSpec AssetKind = "published-spec"
	AssetKindDocument      AssetKind = "document"
	AssetKindCodeChunk     AssetKind = "code-chunk"
	AssetKindIRNode        AssetKind = "ir-node"
	AssetKindArtifact      AssetKind = "artifact"
)

type ArtifactKind string

const (
	ArtifactKindTemplate   ArtifactKind = "template"
	ArtifactKindPolicyPack ArtifactKind = "policy-pack"
	ArtifactKindExportPack ArtifactKind = "export-pack"
	ArtifactKindPromptPack ArtifactKind = "prompt-pack"
)

type PublicationStatus string

const (
	PublicationStatusDraft     PublicationStatus = "draft"
	PublicationStatusPublished PublicationStatus = "published"
	PublicationStatusArchived  PublicationStatus = "archived"
)

type PublishedSpec struct {
	ID          string            `json:"id"`
	Slug        string            `json:"slug"`
	Title       string            `json:"title"`
	Summary     string            `json:"summary"`
	Owner       string            `json:"owner"`
	Version     string            `json:"version"`
	Visibility  string            `json:"visibility"`
	Status      PublicationStatus `json:"status"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	PublishedAt time.Time         `json:"published_at"`
}

type SearchAsset struct {
	ID             string            `json:"id"`
	SpecID         string            `json:"spec_id,omitempty"`
	Kind           AssetKind         `json:"kind"`
	SourceKind     SourceKind        `json:"source_kind"`
	URI            string            `json:"uri"`
	Title          string            `json:"title"`
	Language       string            `json:"language,omitempty"`
	Checksum       string            `json:"checksum,omitempty"`
	Keywords       []string          `json:"keywords,omitempty"`
	EmbeddingRefs  []string          `json:"embedding_refs,omitempty"`
	ChunkIDs       []string          `json:"chunk_ids,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	IndexedAt      time.Time         `json:"indexed_at"`
	LastEnrichedAt time.Time         `json:"last_enriched_at,omitempty"`
}

type ReusableArtifact struct {
	ID             string            `json:"id"`
	SpecID         string            `json:"spec_id,omitempty"`
	Kind           ArtifactKind      `json:"kind"`
	Title          string            `json:"title"`
	Format         string            `json:"format"`
	Location       string            `json:"location"`
	Compatibility  []string          `json:"compatibility,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	PublishedAt    time.Time         `json:"published_at"`
	LastVerifiedAt time.Time         `json:"last_verified_at,omitempty"`
}

type AssetLink struct {
	ID           string            `json:"id"`
	SourceAsset  string            `json:"source_asset"`
	TargetAsset  string            `json:"target_asset"`
	Relation     string            `json:"relation"`
	Confidence   float64           `json:"confidence,omitempty"`
	Provenance   string            `json:"provenance,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	DiscoveredAt time.Time         `json:"discovered_at"`
}

type IngestionRun struct {
	ID          string            `json:"id"`
	Stage       string            `json:"stage"`
	SourceRef   string            `json:"source_ref"`
	Status      string            `json:"status"`
	Diagnostics []string          `json:"diagnostics,omitempty"`
	Metrics     map[string]uint64 `json:"metrics,omitempty"`
	StartedAt   time.Time         `json:"started_at"`
	FinishedAt  time.Time         `json:"finished_at,omitempty"`
}

type RetrievalDocument struct {
	Asset       SearchAsset          `json:"asset"`
	Text        string               `json:"text"`
	SparseTerms []string             `json:"sparse_terms,omitempty"`
	Vectors     map[string][]float32 `json:"vectors,omitempty"`
}

type CorpusStore interface {
	Save(ctx context.Context, document SourceDocument) error
	List(ctx context.Context) ([]SourceDocument, error)
}

type MarketplaceCatalogStore interface {
	SaveSpec(ctx context.Context, spec PublishedSpec) error
	ListSpecs(ctx context.Context) ([]PublishedSpec, error)
	SaveArtifact(ctx context.Context, artifact ReusableArtifact) error
	ListArtifacts(ctx context.Context, specID string) ([]ReusableArtifact, error)
}

type SearchAssetStore interface {
	SaveAsset(ctx context.Context, asset SearchAsset) error
	ListAssets(ctx context.Context, specID string) ([]SearchAsset, error)
	SaveLink(ctx context.Context, link AssetLink) error
	ListLinks(ctx context.Context, assetID string) ([]AssetLink, error)
}

type AnalyticsEventStore interface {
	RecordIngestionRun(ctx context.Context, run IngestionRun) error
	ListIngestionRuns(ctx context.Context, sourceRef string) ([]IngestionRun, error)
}

type CacheStore interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type WorkQueue interface {
	Enqueue(ctx context.Context, queue string, payload []byte) error
	Dequeue(ctx context.Context, queue string) ([]byte, error)
}

type RetrievalIndex interface {
	Upsert(ctx context.Context, documents []RetrievalDocument) error
	Search(ctx context.Context, query string, limit int, filters map[string]string) ([]RetrievalResult, error)
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
