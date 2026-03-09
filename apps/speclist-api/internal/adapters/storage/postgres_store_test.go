package storage

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

func TestPostgresStoreBootstrapsSchema(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(postgresCorpusSchema)).WillReturnResult(sqlmock.NewResult(0, 0))

	if _, err := newPostgresStoreWithDB(db); err != nil {
		t.Fatalf("newPostgresStoreWithDB: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestPostgresStoreSaveReplacesDocumentChunks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(postgresCorpusSchema)).WillReturnResult(sqlmock.NewResult(0, 0))
	store, err := newPostgresStoreWithDB(db)
	if err != nil {
		t.Fatalf("newPostgresStoreWithDB: %v", err)
	}

	document := domain.SourceDocument{
		ID:         "doc-1",
		Kind:       domain.SourceKindSpec,
		Title:      "Platform Spec",
		Location:   "openspec/specs/platform.md",
		ImportedAt: time.Unix(1700000000, 0).UTC(),
		Metadata:   map[string]string{"origin": "repository"},
		Chunks: []domain.Chunk{
			{
				ID:       "chunk-1",
				SourceID: "doc-1",
				Section:  "Goals",
				Text:     "Persist to postgres.",
				Citation: "platform.md > Goals",
				Metadata: map[string]string{"rank": "1"},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO source_documents`).WithArgs(
		document.ID,
		string(document.Kind),
		document.Title,
		document.Location,
		document.ImportedAt,
		[]byte(`{"origin":"repository"}`),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM source_chunks WHERE source_id = \$1`).WithArgs(document.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO source_chunks`).WithArgs(
		"chunk-1",
		document.ID,
		0,
		"Goals",
		"Persist to postgres.",
		"platform.md > Goals",
		[]byte(`{"rank":"1"}`),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := store.Save(context.Background(), document); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestPostgresStoreListReturnsDocumentsWithChunks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(postgresCorpusSchema)).WillReturnResult(sqlmock.NewResult(0, 0))
	store, err := newPostgresStoreWithDB(db)
	if err != nil {
		t.Fatalf("newPostgresStoreWithDB: %v", err)
	}

	importedAt := time.Unix(1700000000, 0).UTC()
	documentRows := sqlmock.NewRows([]string{"id", "kind", "title", "location", "imported_at", "metadata"}).
		AddRow("doc-1", "spec", "Platform Spec", "openspec/specs/platform.md", importedAt, []byte(`{"origin":"repository"}`))
	chunkRows := sqlmock.NewRows([]string{"id", "section", "body_text", "citation", "metadata"}).
		AddRow("chunk-1", "Goals", "Persist to postgres.", "platform.md > Goals", []byte(`{"rank":"1"}`))

	mock.ExpectQuery(`SELECT id, kind, title, location, imported_at, metadata`).WillReturnRows(documentRows)
	mock.ExpectQuery(`SELECT id, section, body_text, citation, metadata`).WithArgs("doc-1").WillReturnRows(chunkRows)

	documents, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(documents))
	}
	if len(documents[0].Chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(documents[0].Chunks))
	}
	if documents[0].Chunks[0].SourceID != "doc-1" {
		t.Fatalf("unexpected source id: %s", documents[0].Chunks[0].SourceID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}
