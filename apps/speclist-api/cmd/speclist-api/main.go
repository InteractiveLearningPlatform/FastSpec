package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/httpapi"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/ingest"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/specs"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/storage"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/app"
)

func main() {
	addr := envOrDefault("SPECLIST_ADDR", ":8080")
	dataFile := envOrDefault("SPECLIST_DATA_FILE", filepath.Join("data", "corpus.json"))
	repoRoot := envOrDefault("SPECLIST_REPO_ROOT", filepath.Join("..", ".."))

	store := storage.NewFileStore(dataFile)
	service := app.NewService(
		store,
		ingest.NewDOCXImporter(),
		ingest.NewConfluenceImporter(nil),
		specs.NewIndexer(),
	)
	server := httpapi.NewServer(service, repoRoot)

	log.Printf("speclist-api listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
