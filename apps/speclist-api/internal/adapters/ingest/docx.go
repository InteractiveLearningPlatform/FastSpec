package ingest

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type DOCXImporter struct{}

func NewDOCXImporter() *DOCXImporter {
	return &DOCXImporter{}
}

func (i *DOCXImporter) Import(_ context.Context, filename string, contents []byte) (domain.SourceDocument, error) {
	if len(contents) == 0 {
		return domain.SourceDocument{}, fmt.Errorf("docx file is empty")
	}

	reader, err := zip.NewReader(bytes.NewReader(contents), int64(len(contents)))
	if err != nil {
		return domain.SourceDocument{}, fmt.Errorf("open docx: %w", err)
	}

	documentXML, err := readZipFile(reader, "word/document.xml")
	if err != nil {
		return domain.SourceDocument{}, err
	}

	text, err := extractDOCXText(documentXML)
	if err != nil {
		return domain.SourceDocument{}, err
	}

	chunks := chunkStructuredText(text, filepath.Base(filename), "docx")
	if len(chunks) == 0 {
		return domain.SourceDocument{}, fmt.Errorf("docx did not contain extractable text")
	}

	sourceID := stableID("docx", filepath.Base(filename))
	for index := range chunks {
		chunks[index].SourceID = sourceID
	}

	title := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	return domain.SourceDocument{
		ID:         sourceID,
		Kind:       domain.SourceKindDocx,
		Title:      title,
		Location:   filename,
		ImportedAt: time.Now().UTC(),
		Metadata: map[string]string{
			"filename": filepath.Base(filename),
		},
		Chunks: chunks,
	}, nil
}

func readZipFile(reader *zip.Reader, target string) ([]byte, error) {
	for _, file := range reader.File {
		if file.Name != target {
			continue
		}
		handle, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer handle.Close()
		return io.ReadAll(handle)
	}
	return nil, fmt.Errorf("docx is missing %s", target)
}

func extractDOCXText(documentXML []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(documentXML))
	lines := make([]string, 0)
	current := strings.Builder{}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("decode docx xml: %w", err)
		}

		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "t" {
				var value string
				if err := decoder.DecodeElement(&value, &element); err != nil {
					return "", fmt.Errorf("decode docx text node: %w", err)
				}
				if current.Len() > 0 {
					current.WriteString(" ")
				}
				current.WriteString(strings.TrimSpace(value))
			}
		case xml.EndElement:
			if element.Name.Local == "p" {
				line := strings.TrimSpace(current.String())
				if line != "" {
					lines = append(lines, line)
				}
				current.Reset()
			}
		}
	}

	return strings.Join(lines, "\n"), nil
}
