package ingest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

func TestConfluenceImporterFetchesPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"id":"123",
			"type":"page",
			"title":"Spec Notes",
			"body":{"storage":{"value":"<h1>Context</h1><p>Grounded retrieval for new specs.</p>"}}
		}`))
	}))
	defer server.Close()

	importer := NewConfluenceImporter(server.Client())
	document, err := importer.Import(context.Background(), domain.ConfluenceImportRequest{
		BaseURL: server.URL,
		PageID:  "123",
	})
	if err != nil {
		t.Fatalf("import confluence page: %v", err)
	}
	if document.Title != "Spec Notes" {
		t.Fatalf("unexpected title: %s", document.Title)
	}
	if len(document.Chunks) == 0 {
		t.Fatal("expected chunks from confluence page")
	}
}
