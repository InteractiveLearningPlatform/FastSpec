package specs

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type Indexer struct{}

func NewIndexer() *Indexer {
	return &Indexer{}
}

func (i *Indexer) Index(_ context.Context, repoRoot string) ([]domain.SourceDocument, error) {
	patterns := []string{
		filepath.Join(repoRoot, "openspec", "specs"),
		filepath.Join(repoRoot, "openspec", "changes"),
		filepath.Join(repoRoot, "examples"),
		filepath.Join(repoRoot, "templates"),
	}

	documents := make([]domain.SourceDocument, 0)
	seen := make(map[string]struct{})
	for _, root := range patterns {
		if _, err := os.Stat(root); err != nil {
			continue
		}

		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
				return nil
			}
			if _, ok := seen[path]; ok {
				return nil
			}
			seen[path] = struct{}{}

			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			relative := path
			if rel, err := filepath.Rel(repoRoot, path); err == nil {
				relative = rel
			}
			title := deriveTitle(relative, string(contents))
			chunks := ingestChunker(string(contents), title)
			sourceID := "spec-" + strings.ReplaceAll(relative, string(filepath.Separator), "-")
			for index := range chunks {
				chunks[index].SourceID = sourceID
				chunks[index].Citation = relative + " > " + chunks[index].Section
			}

			documents = append(documents, domain.SourceDocument{
				ID:         sourceID,
				Kind:       domain.SourceKindSpec,
				Title:      title,
				Location:   relative,
				ImportedAt: time.Now().UTC(),
				Metadata: map[string]string{
					"path": relative,
				},
				Chunks: chunks,
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return documents, nil
}

func deriveTitle(relativePath string, contents string) string {
	for _, line := range strings.Split(contents, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return relativePath
}

func ingestChunker(contents string, title string) []domain.Chunk {
	return ingestChunkStructuredText(contents, title)
}

func ingestChunkStructuredText(contents string, title string) []domain.Chunk {
	lines := strings.Split(contents, "\n")
	chunks := make([]domain.Chunk, 0)
	section := title
	body := make([]string, 0)
	flush := func() {
		text := strings.TrimSpace(strings.Join(body, "\n"))
		if text == "" {
			body = body[:0]
			return
		}
		chunks = append(chunks, domain.Chunk{
			ID:       stableChunkID(title, section, len(chunks)),
			Section:  section,
			Text:     text,
			Metadata: map[string]string{"section": section},
		})
		body = body[:0]
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			flush()
			section = strings.TrimSpace(strings.TrimLeft(line, "#"))
			continue
		}
		body = append(body, line)
		if len(strings.Join(body, "\n")) > 600 {
			flush()
		}
	}
	flush()

	if len(chunks) == 0 {
		chunks = append(chunks, domain.Chunk{
			ID:       stableChunkID(title, title, 0),
			Section:  title,
			Text:     strings.TrimSpace(contents),
			Metadata: map[string]string{"section": title},
		})
	}

	return chunks
}

func stableChunkID(title string, section string, index int) string {
	replacer := strings.NewReplacer(" ", "-", "/", "-", "\\", "-", ".", "-")
	return "chunk-" + replacer.Replace(strings.ToLower(title)) + "-" + replacer.Replace(strings.ToLower(section)) + "-" + strconv.Itoa(index)
}
