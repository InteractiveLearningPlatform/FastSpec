package specs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestIndexerLoadsSpecFiles(t *testing.T) {
	root := t.TempDir()
	specDir := filepath.Join(root, "openspec", "specs", "demo")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte("# Demo Spec\n\n## Requirement\nSystem MUST preserve citations.\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	indexer := NewIndexer()
	documents, err := indexer.Index(context.Background(), root)
	if err != nil {
		t.Fatalf("index specs: %v", err)
	}
	if len(documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(documents))
	}
	if documents[0].Location != filepath.Join("openspec", "specs", "demo", "spec.md") {
		t.Fatalf("unexpected location: %s", documents[0].Location)
	}
}
