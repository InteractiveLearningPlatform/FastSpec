package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type Service struct {
	store              domain.CorpusStore
	docxImporter       domain.DocxImporter
	confluenceImporter domain.ConfluenceImporter
	specIndexer        domain.SpecIndexer
	repoRoot           string
}

func NewService(
	store domain.CorpusStore,
	docxImporter domain.DocxImporter,
	confluenceImporter domain.ConfluenceImporter,
	specIndexer domain.SpecIndexer,
	repoRoot string,
) *Service {
	return &Service{
		store:              store,
		docxImporter:       docxImporter,
		confluenceImporter: confluenceImporter,
		specIndexer:        specIndexer,
		repoRoot:           repoRoot,
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

func (s *Service) ExportDraft(_ context.Context, request domain.DraftExportRequest) (domain.DraftExportResult, error) {
	if strings.TrimSpace(request.Draft.Title) == "" {
		return domain.DraftExportResult{}, fmt.Errorf("draft title is required")
	}
	if len(request.Draft.Sections) == 0 {
		return domain.DraftExportResult{}, fmt.Errorf("draft must contain at least one section")
	}

	format := request.Format
	if format == "" {
		format = domain.ExportFormatOpenSpecMarkdown
	}

	var primaryPath string
	var sidecarPath string
	var primaryContents string
	if request.OpenSpecTarget != nil {
		if format != domain.ExportFormatOpenSpecMarkdown {
			return domain.DraftExportResult{}, fmt.Errorf("openspec change targets require openspec-markdown format")
		}
		var err error
		primaryPath, sidecarPath, err = s.resolveOpenSpecTarget(*request.OpenSpecTarget)
		if err != nil {
			return domain.DraftExportResult{}, err
		}
		primaryContents = renderOpenSpecMarkdown(request.Draft)
	} else {
		if strings.TrimSpace(request.TargetDir) == "" {
			return domain.DraftExportResult{}, fmt.Errorf("target_dir is required")
		}
		if strings.TrimSpace(request.TargetName) == "" {
			return domain.DraftExportResult{}, fmt.Errorf("target_name is required")
		}
		targetDir := filepath.Clean(request.TargetDir)
		targetName := sanitizeFilename(request.TargetName)
		if targetName == "" {
			return domain.DraftExportResult{}, fmt.Errorf("target_name must contain letters or numbers")
		}

		switch format {
		case domain.ExportFormatOpenSpecMarkdown:
			primaryPath = filepath.Join(targetDir, targetName+".md")
			sidecarPath = filepath.Join(targetDir, targetName+".sources.json")
			primaryContents = renderOpenSpecMarkdown(request.Draft)
		case domain.ExportFormatFastSpecYAML:
			primaryPath = filepath.Join(targetDir, targetName+".fastspec.yaml")
			sidecarPath = filepath.Join(targetDir, targetName+".sources.json")
			primaryContents = renderFastSpecYAML(request.Draft, targetName)
		default:
			return domain.DraftExportResult{}, fmt.Errorf("unsupported export format %q", format)
		}
	}

	if err := os.MkdirAll(filepath.Dir(primaryPath), 0o755); err != nil {
		return domain.DraftExportResult{}, err
	}
	if err := ensureFilesDoNotExist(primaryPath, sidecarPath); err != nil {
		return domain.DraftExportResult{}, err
	}

	sidecarContents, err := json.MarshalIndent(renderCitationSidecar(request.Draft, primaryPath), "", "  ")
	if err != nil {
		return domain.DraftExportResult{}, err
	}
	if err := os.WriteFile(primaryPath, []byte(primaryContents), 0o644); err != nil {
		return domain.DraftExportResult{}, err
	}
	if err := os.WriteFile(sidecarPath, sidecarContents, 0o644); err != nil {
		return domain.DraftExportResult{}, err
	}

	return domain.DraftExportResult{
		Format: format,
		Artifacts: []domain.ExportArtifact{
			{Path: primaryPath, Description: "primary exported artifact"},
			{Path: sidecarPath, Description: "citation sidecar"},
		},
	}, nil
}

func (s *Service) ListOpenSpecChanges(_ context.Context) ([]domain.OpenSpecChange, error) {
	if strings.TrimSpace(s.repoRoot) == "" {
		return nil, fmt.Errorf("repo root is not configured")
	}
	changesDir := filepath.Join(s.repoRoot, "openspec", "changes")
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, err
	}

	changes := make([]domain.OpenSpecChange, 0)
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "archive" {
			continue
		}
		changePath := filepath.Join(changesDir, entry.Name())
		if _, err := os.Stat(filepath.Join(changePath, ".openspec.yaml")); err != nil {
			continue
		}
		changes = append(changes, domain.OpenSpecChange{
			Name:      entry.Name(),
			Artifacts: []string{"proposal", "design", "tasks", "spec"},
		})
	}
	slices.SortFunc(changes, func(left, right domain.OpenSpecChange) int {
		return strings.Compare(left.Name, right.Name)
	})
	return changes, nil
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

func renderOpenSpecMarkdown(draft domain.DraftSpec) string {
	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(draft.Title)
	builder.WriteString("\n\n")
	builder.WriteString(draft.Summary)
	builder.WriteString("\n\n")
	builder.WriteString("Query: `")
	builder.WriteString(draft.Query)
	builder.WriteString("`\n\n")
	for _, section := range draft.Sections {
		builder.WriteString("## ")
		builder.WriteString(section.Heading)
		builder.WriteString("\n\n")
		builder.WriteString(section.Body)
		builder.WriteString("\n\n")
		if len(section.Citations) > 0 {
			builder.WriteString("Citations:\n")
			for _, citation := range section.Citations {
				builder.WriteString("- ")
				builder.WriteString(citation)
				builder.WriteString("\n")
			}
			builder.WriteString("\n")
		}
	}
	return strings.TrimSpace(builder.String()) + "\n"
}

func renderFastSpecYAML(draft domain.DraftSpec, targetName string) string {
	var builder strings.Builder
	builder.WriteString("apiVersion: speclist.fastspec.dev/v0alpha1\n")
	builder.WriteString("kind: FastSpecDraft\n")
	builder.WriteString("metadata:\n")
	builder.WriteString("  id: ")
	builder.WriteString(targetName)
	builder.WriteString("\n")
	builder.WriteString("  title: ")
	builder.WriteString(yamlQuote(draft.Title))
	builder.WriteString("\n")
	builder.WriteString("  summary: ")
	builder.WriteString(yamlQuote(draft.Summary))
	builder.WriteString("\n")
	builder.WriteString("spec:\n")
	builder.WriteString("  query: ")
	builder.WriteString(yamlQuote(draft.Query))
	builder.WriteString("\n")
	builder.WriteString("  sourceCount: ")
	builder.WriteString(fmt.Sprintf("%d", draft.SourceCount))
	builder.WriteString("\n")
	builder.WriteString("  sections:\n")
	for _, section := range draft.Sections {
		builder.WriteString("    - heading: ")
		builder.WriteString(yamlQuote(section.Heading))
		builder.WriteString("\n")
		builder.WriteString("      body: |-\n")
		for _, line := range strings.Split(section.Body, "\n") {
			builder.WriteString("        ")
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		builder.WriteString("      citations:\n")
		if len(section.Citations) == 0 {
			builder.WriteString("        []\n")
			continue
		}
		for _, citation := range section.Citations {
			builder.WriteString("        - ")
			builder.WriteString(yamlQuote(citation))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func renderCitationSidecar(draft domain.DraftSpec, primaryPath string) map[string]any {
	sections := make([]map[string]any, 0, len(draft.Sections))
	for _, section := range draft.Sections {
		sections = append(sections, map[string]any{
			"heading":   section.Heading,
			"citations": section.Citations,
		})
	}
	return map[string]any{
		"title":    draft.Title,
		"query":    draft.Query,
		"artifact": primaryPath,
		"sections": sections,
	}
}

func ensureFilesDoNotExist(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("export target already exists: %s", path)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func sanitizeFilename(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	replacer := strings.NewReplacer(" ", "-", "/", "-", "\\", "-", ":", "-", ".", "-", "`", "", "\"", "", "'", "")
	input = replacer.Replace(input)
	input = strings.Trim(input, "-")
	return input
}

func yamlQuote(input string) string {
	input = strings.ReplaceAll(input, "\"", "\\\"")
	return "\"" + input + "\""
}

func (s *Service) resolveOpenSpecTarget(target domain.OpenSpecExportTarget) (string, string, error) {
	if strings.TrimSpace(s.repoRoot) == "" {
		return "", "", fmt.Errorf("repo root is not configured")
	}
	if strings.TrimSpace(target.ChangeName) == "" || strings.TrimSpace(target.Artifact) == "" {
		return "", "", fmt.Errorf("openspec change_name and artifact are required")
	}

	changePath := filepath.Join(s.repoRoot, "openspec", "changes", target.ChangeName)
	if _, err := os.Stat(filepath.Join(changePath, ".openspec.yaml")); err != nil {
		return "", "", fmt.Errorf("openspec change not found: %s", target.ChangeName)
	}

	var primaryPath string
	switch target.Artifact {
	case "proposal":
		primaryPath = filepath.Join(changePath, "proposal.md")
	case "design":
		primaryPath = filepath.Join(changePath, "design.md")
	case "tasks":
		primaryPath = filepath.Join(changePath, "tasks.md")
	case "spec":
		if strings.TrimSpace(target.CapabilityName) == "" {
			return "", "", fmt.Errorf("capability_name is required for spec export")
		}
		capability := sanitizeFilename(target.CapabilityName)
		if capability == "" {
			return "", "", fmt.Errorf("capability_name must contain letters or numbers")
		}
		primaryPath = filepath.Join(changePath, "specs", capability, "spec.md")
	default:
		return "", "", fmt.Errorf("unsupported openspec artifact %q", target.Artifact)
	}

	sidecarPath := strings.TrimSuffix(primaryPath, filepath.Ext(primaryPath)) + ".sources.json"
	return primaryPath, sidecarPath, nil
}
