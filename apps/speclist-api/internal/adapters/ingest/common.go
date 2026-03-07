package ingest

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

func stableID(kind string, value string) string {
	sum := sha1.Sum([]byte(kind + ":" + value))
	return kind + "-" + hex.EncodeToString(sum[:6])
}

func chunkStructuredText(text string, title string, sourceType string) []domain.Chunk {
	lines := strings.Split(text, "\n")
	chunks := make([]domain.Chunk, 0)
	section := title
	buffer := make([]string, 0)
	flush := func() {
		if len(buffer) == 0 {
			return
		}
		body := strings.TrimSpace(strings.Join(buffer, " "))
		if body == "" {
			buffer = buffer[:0]
			return
		}
		chunks = append(chunks, domain.Chunk{
			ID:       stableID(sourceType+"-chunk", fmt.Sprintf("%s:%d:%s", title, len(chunks), section)),
			Section:  section,
			Text:     body,
			Citation: fmt.Sprintf("%s > %s", title, section),
			Metadata: map[string]string{"section": section},
		})
		buffer = buffer[:0]
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			flush()
			continue
		}

		if looksLikeHeading(line) {
			flush()
			section = line
			continue
		}

		buffer = append(buffer, line)
		if len(strings.Join(buffer, " ")) > 420 {
			flush()
		}
	}
	flush()

	if len(chunks) == 0 && strings.TrimSpace(text) != "" {
		chunks = append(chunks, domain.Chunk{
			ID:       stableID(sourceType+"-chunk", title),
			Section:  title,
			Text:     strings.TrimSpace(text),
			Citation: title,
			Metadata: map[string]string{"section": title},
		})
	}

	return chunks
}

func looksLikeHeading(line string) bool {
	return strings.HasPrefix(line, "#") ||
		strings.HasSuffix(line, ":") ||
		(strings.Count(line, " ") <= 6 && line == strings.ToUpper(line))
}
