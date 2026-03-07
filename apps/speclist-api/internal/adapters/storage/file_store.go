package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type FileStore struct {
	path string
	mu   sync.Mutex
}

type persistedCorpus struct {
	Documents []domain.SourceDocument `json:"documents"`
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Save(_ context.Context, document domain.SourceDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	corpus, err := s.load()
	if err != nil {
		return err
	}

	replaced := false
	for index, existing := range corpus.Documents {
		if existing.ID == document.ID {
			corpus.Documents[index] = document
			replaced = true
			break
		}
	}
	if !replaced {
		corpus.Documents = append(corpus.Documents, document)
	}

	return s.persist(corpus)
}

func (s *FileStore) List(_ context.Context) ([]domain.SourceDocument, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	corpus, err := s.load()
	if err != nil {
		return nil, err
	}
	return corpus.Documents, nil
}

func (s *FileStore) load() (persistedCorpus, error) {
	if s.path == "" {
		return persistedCorpus{}, errors.New("file store path is required")
	}

	contents, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return persistedCorpus{}, nil
		}
		return persistedCorpus{}, err
	}

	var corpus persistedCorpus
	if err := json.Unmarshal(contents, &corpus); err != nil {
		return persistedCorpus{}, err
	}
	return corpus, nil
}

func (s *FileStore) persist(corpus persistedCorpus) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	contents, err := json.MarshalIndent(corpus, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, contents, 0o644)
}
