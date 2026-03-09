package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

const postgresCorpusSchema = `
CREATE TABLE IF NOT EXISTS source_documents (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    title TEXT NOT NULL,
    location TEXT NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS source_chunks (
    id TEXT PRIMARY KEY,
    source_id TEXT NOT NULL REFERENCES source_documents(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    section TEXT NOT NULL,
    body_text TEXT NOT NULL,
    citation TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS source_chunks_source_position_idx ON source_chunks (source_id, position);
`

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return nil, fmt.Errorf("postgres dsn is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	store, err := newPostgresStoreWithDB(db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func newPostgresStoreWithDB(db *sql.DB) (*PostgresStore, error) {
	store := &PostgresStore{db: db}
	if err := store.bootstrap(context.Background()); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *PostgresStore) Save(ctx context.Context, document domain.SourceDocument) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	documentMetadata, err := marshalMetadata(document.Metadata)
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(
		ctx,
		`INSERT INTO source_documents (id, kind, title, location, imported_at, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (id) DO UPDATE
		 SET kind = EXCLUDED.kind,
		     title = EXCLUDED.title,
		     location = EXCLUDED.location,
		     imported_at = EXCLUDED.imported_at,
		     metadata = EXCLUDED.metadata`,
		document.ID,
		string(document.Kind),
		document.Title,
		document.Location,
		document.ImportedAt,
		documentMetadata,
	); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM source_chunks WHERE source_id = $1`, document.ID); err != nil {
		return err
	}

	for index, chunk := range document.Chunks {
		chunkMetadata, chunkErr := marshalMetadata(chunk.Metadata)
		if chunkErr != nil {
			err = chunkErr
			return err
		}

		if _, err = tx.ExecContext(
			ctx,
			`INSERT INTO source_chunks (id, source_id, position, section, body_text, citation, metadata)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			chunk.ID,
			document.ID,
			index,
			chunk.Section,
			chunk.Text,
			chunk.Citation,
			chunkMetadata,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresStore) List(ctx context.Context) ([]domain.SourceDocument, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, kind, title, location, imported_at, metadata
		 FROM source_documents
		 ORDER BY imported_at ASC, id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]domain.SourceDocument, 0)
	for rows.Next() {
		var (
			document     domain.SourceDocument
			kind         string
			metadataJSON []byte
		)

		if err := rows.Scan(&document.ID, &kind, &document.Title, &document.Location, &document.ImportedAt, &metadataJSON); err != nil {
			return nil, err
		}
		document.Kind = domain.SourceKind(kind)
		document.Metadata, err = unmarshalMetadata(metadataJSON)
		if err != nil {
			return nil, err
		}

		document.Chunks, err = s.listChunks(ctx, document.ID)
		if err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return documents, nil
}

func (s *PostgresStore) bootstrap(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, postgresCorpusSchema)
	return err
}

func (s *PostgresStore) listChunks(ctx context.Context, sourceID string) ([]domain.Chunk, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, section, body_text, citation, metadata
		 FROM source_chunks
		 WHERE source_id = $1
		 ORDER BY position ASC`,
		sourceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chunks := make([]domain.Chunk, 0)
	for rows.Next() {
		var (
			chunk        domain.Chunk
			metadataJSON []byte
		)

		if err := rows.Scan(&chunk.ID, &chunk.Section, &chunk.Text, &chunk.Citation, &metadataJSON); err != nil {
			return nil, err
		}
		chunk.SourceID = sourceID
		chunk.Metadata, err = unmarshalMetadata(metadataJSON)
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return chunks, nil
}

func marshalMetadata(metadata map[string]string) ([]byte, error) {
	if len(metadata) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(metadata)
}

func unmarshalMetadata(contents []byte) (map[string]string, error) {
	if len(contents) == 0 || string(contents) == "null" {
		return nil, nil
	}

	var metadata map[string]string
	if err := json.Unmarshal(contents, &metadata); err != nil {
		return nil, err
	}
	if len(metadata) == 0 {
		return nil, nil
	}
	return metadata, nil
}
