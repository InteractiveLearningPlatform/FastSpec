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

func (s *Service) InspectSource(ctx context.Context, sourceID string) (domain.SourceDocument, error) {
	sourceID = strings.TrimSpace(sourceID)
	if sourceID == "" {
		return domain.SourceDocument{}, fmt.Errorf("source_id is required")
	}

	documents, err := s.store.List(ctx)
	if err != nil {
		return domain.SourceDocument{}, err
	}
	for _, document := range documents {
		if document.ID == sourceID {
			return document, nil
		}
	}

	return domain.SourceDocument{}, fmt.Errorf("source %q was not found", sourceID)
}

func (s *Service) InspectCitation(ctx context.Context, citation string) (domain.CitationInspection, error) {
	citation = strings.TrimSpace(citation)
	if citation == "" {
		return domain.CitationInspection{}, fmt.Errorf("citation is required")
	}

	documents, err := s.store.List(ctx)
	if err != nil {
		return domain.CitationInspection{}, err
	}

	for _, document := range documents {
		source := domain.SourceStub{
			ID:       document.ID,
			Kind:     document.Kind,
			Title:    document.Title,
			Location: document.Location,
			Metadata: document.Metadata,
		}
		for _, chunk := range document.Chunks {
			if chunk.Citation != citation {
				continue
			}
			return domain.CitationInspection{
				Citation: citation,
				Source:   source,
				Chunk:    chunk,
			}, nil
		}
	}

	return domain.CitationInspection{}, fmt.Errorf("citation %q was not found", citation)
}

func (s *Service) Search(ctx context.Context, query string, limit int, filters domain.RetrievalFilter) (domain.RetrievalBundle, error) {
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
		if !matchesFilters(document, filters) {
			continue
		}
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

	return domain.RetrievalBundle{Query: query, Filters: normalizeFilters(filters), Results: results}, nil
}

func (s *Service) DraftSpec(ctx context.Context, query string, title string, format string, limit int, filters domain.RetrievalFilter) (domain.DraftSpec, error) {
	bundle, err := s.Search(ctx, query, limit, filters)
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
		Filters:     bundle.Filters,
		Format:      format,
		Summary:     fmt.Sprintf("Drafted from %d grounded retrieval result(s).", len(bundle.Results)),
		Sections:    sections,
		SourceCount: len(bundle.Results),
	}, nil
}

func (s *Service) ExportDraft(_ context.Context, request domain.DraftExportRequest) (domain.DraftExportResult, error) {
	draft, err := normalizeDraft(request.Draft)
	if err != nil {
		return domain.DraftExportResult{}, err
	}

	if strings.TrimSpace(draft.Title) == "" {
		return domain.DraftExportResult{}, fmt.Errorf("draft title is required")
	}
	if len(draft.Sections) == 0 {
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
		primaryPath, sidecarPath, err = s.resolveOpenSpecTarget(*request.OpenSpecTarget)
		if err != nil {
			return domain.DraftExportResult{}, err
		}
		primaryContents = renderOpenSpecArtifact(draft, *request.OpenSpecTarget)
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
			primaryContents = renderOpenSpecMarkdown(draft)
		case domain.ExportFormatFastSpecYAML:
			primaryPath = filepath.Join(targetDir, targetName+".fastspec.yaml")
			sidecarPath = filepath.Join(targetDir, targetName+".sources.json")
			primaryContents = renderFastSpecYAML(draft, targetName)
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

	sidecarContents, err := json.MarshalIndent(renderCitationSidecar(draft, primaryPath), "", "  ")
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

func normalizeDraft(draft domain.DraftSpec) (domain.DraftSpec, error) {
	draft.Title = strings.TrimSpace(draft.Title)
	draft.Query = strings.TrimSpace(draft.Query)
	draft.Summary = strings.TrimSpace(draft.Summary)
	if draft.Title == "" {
		return domain.DraftSpec{}, fmt.Errorf("draft title is required")
	}

	sections := make([]domain.DraftSection, 0, len(draft.Sections))
	for _, section := range draft.Sections {
		normalized := domain.DraftSection{
			Heading: strings.TrimSpace(section.Heading),
			Body:    strings.TrimSpace(section.Body),
		}
		if normalized.Heading == "" {
			return domain.DraftSpec{}, fmt.Errorf("draft section heading is required")
		}
		if normalized.Body == "" {
			return domain.DraftSpec{}, fmt.Errorf("draft section body is required")
		}
		for _, citation := range section.Citations {
			citation = strings.TrimSpace(citation)
			if citation == "" {
				continue
			}
			normalized.Citations = append(normalized.Citations, citation)
		}
		sections = append(sections, normalized)
	}
	draft.Sections = sections
	return draft, nil
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

func matchesFilters(document domain.SourceDocument, filters domain.RetrievalFilter) bool {
	normalized := normalizeFilters(filters)
	if len(normalized.Kinds) > 0 && !slices.Contains(normalized.Kinds, document.Kind) {
		return false
	}
	if normalized.Origin != "" && sourceOrigin(document) != normalized.Origin {
		return false
	}
	if normalized.LocationContains != "" && !strings.Contains(strings.ToLower(document.Location), normalized.LocationContains) {
		return false
	}
	return true
}

func normalizeFilters(filters domain.RetrievalFilter) domain.RetrievalFilter {
	normalized := domain.RetrievalFilter{
		Origin:           filters.Origin,
		LocationContains: strings.ToLower(strings.TrimSpace(filters.LocationContains)),
	}
	if len(filters.Kinds) > 0 {
		seen := make(map[domain.SourceKind]struct{}, len(filters.Kinds))
		for _, kind := range filters.Kinds {
			if kind == "" {
				continue
			}
			if _, ok := seen[kind]; ok {
				continue
			}
			seen[kind] = struct{}{}
			normalized.Kinds = append(normalized.Kinds, kind)
		}
		slices.Sort(normalized.Kinds)
	}
	return normalized
}

func sourceOrigin(document domain.SourceDocument) domain.SourceOrigin {
	if document.Kind == domain.SourceKindSpec {
		return domain.SourceOriginRepository
	}
	return domain.SourceOriginImported
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

func renderOpenSpecArtifact(draft domain.DraftSpec, target domain.OpenSpecExportTarget) string {
	switch target.Artifact {
	case "proposal":
		return renderProposalTemplate(draft)
	case "design":
		return renderDesignTemplate(draft)
	case "tasks":
		return renderTasksTemplate(draft)
	case "spec":
		return renderSpecTemplate(draft, target.CapabilityName)
	default:
		return renderOpenSpecMarkdown(draft)
	}
}

func renderProposalTemplate(draft domain.DraftSpec) string {
	return strings.TrimSpace(fmt.Sprintf(`## Why

%s

## What Changes

%s

## Capabilities

### New Capabilities
- %s

### Modified Capabilities
- none yet

## Impact

%s
`,
		findSectionBody(draft, "Why", draft.Summary),
		findSectionBody(draft, "Proposed Requirements", draft.Summary),
		strings.ToLower(strings.ReplaceAll(draft.Title, " ", "-")),
		findSectionBody(draft, "Context", draft.Summary),
	)) + "\n"
}

func renderDesignTemplate(draft domain.DraftSpec) string {
	return strings.TrimSpace(fmt.Sprintf(`## Context

%s

## Goals / Non-Goals

**Goals:**
%s

**Non-Goals:**
- refine after review

## Decisions

%s

## Risks / Trade-offs

- review exported draft assumptions against cited sources
`,
		findSectionBody(draft, "Context", draft.Summary),
		findSectionBody(draft, "Why", draft.Summary),
		findSectionBody(draft, "Proposed Requirements", draft.Summary),
	)) + "\n"
}

func renderTasksTemplate(draft domain.DraftSpec) string {
	lines := bulletize(findSectionBody(draft, "Proposed Requirements", draft.Summary))
	if len(lines) == 0 {
		lines = []string{"refine exported draft into implementation tasks"}
	}

	var builder strings.Builder
	builder.WriteString("## 1. Exported Tasks\n\n")
	for index, line := range lines {
		builder.WriteString(fmt.Sprintf("- [ ] 1.%d %s\n", index+1, sanitizeTaskLine(line)))
	}
	return builder.String()
}

func renderSpecTemplate(draft domain.DraftSpec, capabilityName string) string {
	requirementName := capabilityName
	if strings.TrimSpace(requirementName) == "" {
		requirementName = sanitizeFilename(draft.Title)
	}
	return strings.TrimSpace(fmt.Sprintf(`## ADDED Requirements

### Requirement: %s
%s

#### Scenario: Drafted Scenario
- **WHEN** a contributor applies the exported draft
- **THEN** the resulting spec starts from the grounded requirement content below
- **AND** the artifact remains ready for further refinement

%s
`,
		titleizeRequirement(requirementName),
		findSectionBody(draft, "Why", draft.Summary),
		indentLines(findSectionBody(draft, "Proposed Requirements", draft.Summary)),
	)) + "\n"
}

func renderFastSpecYAML(draft domain.DraftSpec, targetName string) string {
	var builder strings.Builder
	builder.WriteString("apiVersion: speclist.fastspec.dev/v0alpha1\n")
	builder.WriteString("kind: SpecDocumentDraft\n")
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
	builder.WriteString("  draftType: \"structured-spec\"\n")
	builder.WriteString("spec:\n")
	builder.WriteString("  query: ")
	builder.WriteString(yamlQuote(draft.Query))
	builder.WriteString("\n")
	builder.WriteString("  sourceCount: ")
	builder.WriteString(fmt.Sprintf("%d", draft.SourceCount))
	builder.WriteString("\n")
	builder.WriteString("  rationale: |-\n")
	for _, line := range strings.Split(findSectionBody(draft, "Why", draft.Summary), "\n") {
		builder.WriteString("    ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("  context: |-\n")
	for _, line := range strings.Split(findSectionBody(draft, "Context", draft.Summary), "\n") {
		builder.WriteString("    ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("  requirements:\n")
	requirements := bulletize(findSectionBody(draft, "Proposed Requirements", draft.Summary))
	if len(requirements) == 0 {
		requirements = []string{draft.Summary}
	}
	for index, requirement := range requirements {
		builder.WriteString("    - id: ")
		builder.WriteString(yamlQuote(fmt.Sprintf("%s-%d", targetName, index+1)))
		builder.WriteString("\n")
		builder.WriteString("      statement: ")
		builder.WriteString(yamlQuote(requirement))
		builder.WriteString("\n")
	}
	builder.WriteString("  sections:\n")
	for _, section := range draft.Sections {
		builder.WriteString("    - heading: ")
		builder.WriteString(yamlQuote(section.Heading))
		builder.WriteString("\n")
		builder.WriteString("      summary: ")
		builder.WriteString(yamlQuote(singleLine(section.Body)))
		builder.WriteString("\n")
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

func findSectionBody(draft domain.DraftSpec, heading string, fallback string) string {
	for _, section := range draft.Sections {
		if strings.EqualFold(section.Heading, heading) {
			return section.Body
		}
	}
	return fallback
}

func bulletize(input string) []string {
	lines := make([]string, 0)
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "-"))
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func sanitizeTaskLine(input string) string {
	input = strings.TrimSpace(strings.TrimPrefix(input, "MUST "))
	input = strings.TrimSpace(strings.TrimPrefix(input, "reflect:"))
	input = strings.TrimSpace(strings.TrimPrefix(input, "MUST reflect:"))
	if input == "" {
		return "refine exported draft task"
	}
	return strings.ToLower(input[:1]) + input[1:]
}

func titleizeRequirement(input string) string {
	input = strings.ReplaceAll(input, "-", " ")
	input = strings.TrimSpace(input)
	if input == "" {
		return "Exported Draft Requirement"
	}
	return strings.ToUpper(input[:1]) + input[1:]
}

func indentLines(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for index, line := range lines {
		lines[index] = line
	}
	return strings.Join(lines, "\n")
}

func singleLine(input string) string {
	input = strings.TrimSpace(strings.ReplaceAll(input, "\n", " "))
	input = strings.Join(strings.Fields(input), " ")
	return input
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
