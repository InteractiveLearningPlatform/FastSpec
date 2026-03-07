package app

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type Service struct {
	store              domain.CorpusStore
	docxImporter       domain.DocxImporter
	confluenceImporter domain.ConfluenceImporter
	specIndexer        domain.SpecIndexer
}

func NewService(
	store domain.CorpusStore,
	docxImporter domain.DocxImporter,
	confluenceImporter domain.ConfluenceImporter,
	specIndexer domain.SpecIndexer,
) *Service {
	return &Service{
		store:              store,
		docxImporter:       docxImporter,
		confluenceImporter: confluenceImporter,
		specIndexer:        specIndexer,
	}
}

func (s *Service) ImportDOCX(ctx context.Context, filename string, contents []byte) (domain.SourceDocument, error) {
	document, err := s.docxImporter.Import(ctx, filename, contents)
	if err != nil {
		return domain.SourceDocument{}, err
	}

	if err := s.store.Save(ctx, document); err != nil {
		return domain.SourceDocument{}, err
	}

	return document, nil
}

func (s *Service) ImportConfluence(ctx context.Context, request domain.ConfluenceImportRequest) (domain.SourceDocument, error) {
	document, err := s.confluenceImporter.Import(ctx, request)
	if err != nil {
		return domain.SourceDocument{}, err
	}

	if err := s.store.Save(ctx, document); err != nil {
		return domain.SourceDocument{}, err
	}

	return document, nil
}

func (s *Service) IndexSpecs(ctx context.Context, repoRoot string) ([]domain.SourceDocument, error) {
	documents, err := s.specIndexer.Index(ctx, repoRoot)
	if err != nil {
		return nil, err
	}

	for _, document := range documents {
		if err := s.store.Save(ctx, document); err != nil {
			return nil, err
		}
	}

	return documents, nil
}

func (s *Service) ListSources(ctx context.Context) ([]domain.SourceDocument, error) {
	return s.store.List(ctx)
}

func (s *Service) Search(ctx context.Context, query string, limit int) (domain.RetrievalBundle, error) {
	if strings.TrimSpace(query) == "" {
		return domain.RetrievalBundle{}, fmt.Errorf("query is required")
	}

	documents, err := s.store.List(ctx)
	if err != nil {
		return domain.RetrievalBundle{}, err
	}

	if limit <= 0 {
		limit = 8
	}

	terms := tokenize(query)
	results := make([]domain.RetrievalResult, 0)
	for _, document := range documents {
		source := domain.SourceStub{
			ID:       document.ID,
			Kind:     document.Kind,
			Title:    document.Title,
			Location: document.Location,
			Metadata: document.Metadata,
		}

		for _, chunk := range document.Chunks {
			score := scoreChunk(terms, document, chunk)
			if score == 0 {
				continue
			}
			results = append(results, domain.RetrievalResult{
				Chunk:  chunk,
				Source: source,
				Score:  score,
			})
		}
	}

	slices.SortFunc(results, func(left, right domain.RetrievalResult) int {
		if left.Score != right.Score {
			return right.Score - left.Score
		}
		if left.Source.Title != right.Source.Title {
			return strings.Compare(left.Source.Title, right.Source.Title)
		}
		return strings.Compare(left.Chunk.ID, right.Chunk.ID)
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return domain.RetrievalBundle{Query: query, Results: results}, nil
}

func (s *Service) DraftSpec(ctx context.Context, query string, title string, format string, limit int) (domain.DraftSpec, error) {
	bundle, err := s.Search(ctx, query, limit)
	if err != nil {
		return domain.DraftSpec{}, err
	}

	if strings.TrimSpace(title) == "" {
		title = "Generated Spec Draft"
	}
	if strings.TrimSpace(format) == "" {
		format = "openspec-markdown"
	}

	sections := []domain.DraftSection{
		{
			Heading:   "Why",
			Body:      summarizeBundle(bundle.Results, 0, 2),
			Citations: collectCitations(bundle.Results, 0, 2),
		},
		{
			Heading:   "Context",
			Body:      summarizeBundle(bundle.Results, 0, min(4, len(bundle.Results))),
			Citations: collectCitations(bundle.Results, 0, min(4, len(bundle.Results))),
		},
		{
			Heading:   "Proposed Requirements",
			Body:      summarizeRequirements(bundle.Results),
			Citations: collectCitations(bundle.Results, 0, min(6, len(bundle.Results))),
		},
	}

	return domain.DraftSpec{
		Title:       title,
		Query:       query,
		Format:      format,
		Summary:     fmt.Sprintf("Drafted from %d grounded retrieval result(s).", len(bundle.Results)),
		Sections:    sections,
		SourceCount: len(bundle.Results),
	}, nil
}

func summarizeBundle(results []domain.RetrievalResult, start int, end int) string {
	if len(results) == 0 || start >= end || start >= len(results) {
		return "No grounded sources were retrieved."
	}
	if end > len(results) {
		end = len(results)
	}

	lines := make([]string, 0, end-start)
	for _, result := range results[start:end] {
		lines = append(lines, fmt.Sprintf("- %s (%s): %s", result.Source.Title, result.Chunk.Citation, clipText(result.Chunk.Text, 180)))
	}
	return strings.Join(lines, "\n")
}

func summarizeRequirements(results []domain.RetrievalResult) string {
	if len(results) == 0 {
		return "- Add requirements after importing relevant sources."
	}

	lines := make([]string, 0, min(5, len(results)))
	for _, result := range results[:min(5, len(results))] {
		lines = append(lines, fmt.Sprintf("- MUST reflect: %s", clipText(result.Chunk.Text, 160)))
	}
	return strings.Join(lines, "\n")
}

func collectCitations(results []domain.RetrievalResult, start int, end int) []string {
	if len(results) == 0 || start >= end || start >= len(results) {
		return nil
	}
	if end > len(results) {
		end = len(results)
	}

	citations := make([]string, 0, end-start)
	seen := make(map[string]struct{})
	for _, result := range results[start:end] {
		if _, ok := seen[result.Chunk.Citation]; ok {
			continue
		}
		seen[result.Chunk.Citation] = struct{}{}
		citations = append(citations, result.Chunk.Citation)
	}
	return citations
}

func scoreChunk(terms []string, document domain.SourceDocument, chunk domain.Chunk) int {
	haystack := strings.ToLower(strings.Join([]string{
		document.Title,
		document.Location,
		chunk.Section,
		chunk.Text,
		metadataString(document.Metadata),
		metadataString(chunk.Metadata),
	}, " "))

	score := 0
	for _, term := range terms {
		if strings.Contains(haystack, term) {
			score++
		}
	}
	return score
}

func tokenize(input string) []string {
	fields := strings.Fields(strings.ToLower(input))
	terms := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.Trim(field, ".,:;()[]{}\"'")
		if field != "" {
			terms = append(terms, field)
		}
	}
	return terms
}

func metadataString(values map[string]string) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for key, value := range values {
		parts = append(parts, key+" "+value)
	}
	return strings.Join(parts, " ")
}

func clipText(input string, limit int) string {
	input = strings.TrimSpace(input)
	if len(input) <= limit {
		return input
	}
	return strings.TrimSpace(input[:limit-3]) + "..."
}

func min(left int, right int) int {
	if left < right {
		return left
	}
	return right
}
