package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/httpapi"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/ingest"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/specs"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/adapters/storage"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/app"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

func main() {
	addr := envOrDefault("SPECLIST_ADDR", ":8080")
	dataFile := filepath.Join("data", "corpus.json")
	repoRoot := envOrDefault("SPECLIST_REPO_ROOT", filepath.Join("..", ".."))

	store, err := selectStore(os.Getenv, storeFactories{
		newFile: func(path string) domain.CorpusStore {
			return storage.NewFileStore(path)
		},
		newPostgres: func(dsn string) (domain.CorpusStore, error) {
			return storage.NewPostgresStore(dsn)
		},
	}, dataFile)
	if err != nil {
		log.Fatal(err)
	}
	service := app.NewService(
		store,
		ingest.NewDOCXImporter(),
		ingest.NewConfluenceImporter(nil),
		specs.NewIndexer(),
		repoRoot,
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

type storeFactories struct {
	newFile     func(path string) domain.CorpusStore
	newPostgres func(dsn string) (domain.CorpusStore, error)
}

func selectStore(getenv func(string) string, factories storeFactories, defaultDataFile string) (domain.CorpusStore, error) {
	storeKind := strings.TrimSpace(getenv("SPECLIST_STORE_KIND"))
	if storeKind == "" {
		storeKind = "file"
	}

	switch storeKind {
	case "file":
		dataFile := getenv("SPECLIST_DATA_FILE")
		if strings.TrimSpace(dataFile) == "" {
			dataFile = defaultDataFile
		}
		return factories.newFile(dataFile), nil
	case "postgres":
		dsn := strings.TrimSpace(getenv("SPECLIST_POSTGRES_DSN"))
		if dsn == "" {
			return nil, fmt.Errorf("SPECLIST_POSTGRES_DSN is required when SPECLIST_STORE_KIND=postgres")
		}
		return factories.newPostgres(dsn)
	default:
		return nil, fmt.Errorf("unsupported SPECLIST_STORE_KIND %q", storeKind)
	}
}
