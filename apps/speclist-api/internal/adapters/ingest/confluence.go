package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type ConfluenceImporter struct {
	client *http.Client
}

func NewConfluenceImporter(client *http.Client) *ConfluenceImporter {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &ConfluenceImporter{client: client}
}

func (i *ConfluenceImporter) Import(ctx context.Context, request domain.ConfluenceImportRequest) (domain.SourceDocument, error) {
	if strings.TrimSpace(request.BaseURL) == "" || strings.TrimSpace(request.PageID) == "" {
		return domain.SourceDocument{}, fmt.Errorf("base_url and page_id are required")
	}

	url := strings.TrimRight(request.BaseURL, "/") + "/rest/api/content/" + request.PageID + "?expand=body.storage"
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.SourceDocument{}, err
	}
	if request.Token != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+request.Token)
	}
	httpRequest.Header.Set("Accept", "application/json")

	response, err := i.client.Do(httpRequest)
	if err != nil {
		return domain.SourceDocument{}, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return domain.SourceDocument{}, fmt.Errorf("confluence request failed: %s: %s", response.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Title string `json:"title"`
		Body  struct {
			Storage struct {
				Value string `json:"value"`
			} `json:"storage"`
		} `json:"body"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return domain.SourceDocument{}, err
	}

	text := stripHTML(payload.Body.Storage.Value)
	chunks := chunkStructuredText(text, payload.Title, "confluence")
	sourceID := stableID("confluence", payload.ID)
	for index := range chunks {
		chunks[index].SourceID = sourceID
	}

	return domain.SourceDocument{
		ID:         sourceID,
		Kind:       domain.SourceKindConfluence,
		Title:      payload.Title,
		Location:   url,
		ImportedAt: time.Now().UTC(),
		Metadata: map[string]string{
			"page_id": payload.ID,
			"type":    payload.Type,
		},
		Chunks: chunks,
	}, nil
}

var htmlTagPattern = regexp.MustCompile("<[^>]+>")

func stripHTML(input string) string {
	replaced := htmlTagPattern.ReplaceAllString(input, "\n")
	replaced = strings.ReplaceAll(replaced, "&nbsp;", " ")
	replaced = strings.ReplaceAll(replaced, "&amp;", "&")
	lines := strings.Split(replaced, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}
