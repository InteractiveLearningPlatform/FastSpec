package main

import (
	"context"
	"strings"
	"testing"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type stubStore struct{}

func (stubStore) Save(context.Context, domain.SourceDocument) error { return nil }
func (stubStore) List(context.Context) ([]domain.SourceDocument, error) {
	return nil, nil
}

func TestSelectStoreUsesFileStoreByDefault(t *testing.T) {
	var filePath string
	store, err := selectStore(func(string) string { return "" }, storeFactories{
		newFile: func(path string) domain.CorpusStore {
			filePath = path
			return stubStore{}
		},
		newPostgres: func(string) (domain.CorpusStore, error) {
			t.Fatal("postgres factory should not be called")
			return nil, nil
		},
	}, "data/corpus.json")
	if err != nil {
		t.Fatalf("selectStore: %v", err)
	}
	if store == nil {
		t.Fatal("expected store")
	}
	if filePath != "data/corpus.json" {
		t.Fatalf("unexpected file path: %s", filePath)
	}
}

func TestSelectStoreRequiresPostgresDSN(t *testing.T) {
	_, err := selectStore(func(key string) string {
		if key == "SPECLIST_STORE_KIND" {
			return "postgres"
		}
		return ""
	}, storeFactories{
		newFile: func(string) domain.CorpusStore { return stubStore{} },
		newPostgres: func(string) (domain.CorpusStore, error) {
			t.Fatal("postgres factory should not be called")
			return nil, nil
		},
	}, "data/corpus.json")
	if err == nil || !strings.Contains(err.Error(), "SPECLIST_POSTGRES_DSN") {
		t.Fatalf("expected missing dsn error, got %v", err)
	}
}

func TestSelectStoreUsesPostgresFactory(t *testing.T) {
	var receivedDSN string
	store, err := selectStore(func(key string) string {
		switch key {
		case "SPECLIST_STORE_KIND":
			return "postgres"
		case "SPECLIST_POSTGRES_DSN":
			return "postgres://speclist:secret@db/speclist"
		default:
			return ""
		}
	}, storeFactories{
		newFile: func(string) domain.CorpusStore {
			t.Fatal("file factory should not be called")
			return nil
		},
		newPostgres: func(dsn string) (domain.CorpusStore, error) {
			receivedDSN = dsn
			return stubStore{}, nil
		},
	}, "data/corpus.json")
	if err != nil {
		t.Fatalf("selectStore: %v", err)
	}
	if store == nil {
		t.Fatal("expected store")
	}
	if receivedDSN != "postgres://speclist:secret@db/speclist" {
		t.Fatalf("unexpected dsn: %s", receivedDSN)
	}
}
