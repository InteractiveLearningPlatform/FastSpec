package ingest

import (
	"archive/zip"
	"bytes"
	"context"
	"testing"
)

func TestDOCXImporterExtractsParagraphs(t *testing.T) {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)

	file, err := writer.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create zip file: %v", err)
	}
	_, err = file.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<document>
  <body>
    <p><r><t>Overview</t></r></p>
    <p><r><t>System MUST support spec ingestion.</t></r></p>
    <p><r><t>Goals:</t></r></p>
    <p><r><t>Generate drafts with citations.</t></r></p>
  </body>
</document>`))
	if err != nil {
		t.Fatalf("write xml: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	importer := NewDOCXImporter()
	document, err := importer.Import(context.Background(), "notes.docx", buffer.Bytes())
	if err != nil {
		t.Fatalf("import docx: %v", err)
	}
	if document.Kind != "docx" {
		t.Fatalf("unexpected kind: %s", document.Kind)
	}
	if len(document.Chunks) == 0 {
		t.Fatal("expected extracted chunks")
	}
}
