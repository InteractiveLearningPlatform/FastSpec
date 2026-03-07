package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/app"
	"github.com/InteractiveLearningPlatform/FastSpec/apps/speclist-api/internal/domain"
)

type Server struct {
	service  *app.Service
	repoRoot string
}

func NewServer(service *app.Service, repoRoot string) *Server {
	return &Server{service: service, repoRoot: repoRoot}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", s.handleHealth)
	mux.HandleFunc("/api/v1/sources", s.handleListSources)
	mux.HandleFunc("/api/v1/import/docx", s.handleImportDOCX)
	mux.HandleFunc("/api/v1/import/confluence", s.handleImportConfluence)
	mux.HandleFunc("/api/v1/index/specs", s.handleIndexSpecs)
	mux.HandleFunc("/api/v1/search", s.handleSearch)
	mux.HandleFunc("/api/v1/drafts", s.handleDraft)
	return withCORS(mux)
}

func (s *Server) handleHealth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeMethodNotAllowed(writer)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListSources(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeMethodNotAllowed(writer)
		return
	}
	sources, err := s.service.ListSources(request.Context())
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"sources": sources})
}

func (s *Server) handleImportDOCX(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeMethodNotAllowed(writer)
		return
	}
	if err := request.ParseMultipartForm(10 << 20); err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		writeError(writer, http.StatusBadRequest, fmt.Errorf("file is required"))
		return
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}

	document, err := s.service.ImportDOCX(request.Context(), header.Filename, contents)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{"document": document})
}

func (s *Server) handleImportConfluence(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeMethodNotAllowed(writer)
		return
	}
	var payload domain.ConfluenceImportRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	document, err := s.service.ImportConfluence(request.Context(), payload)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{"document": document})
}

func (s *Server) handleIndexSpecs(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeMethodNotAllowed(writer)
		return
	}

	var payload struct {
		RepoRoot string `json:"repo_root"`
	}
	if request.Body != nil {
		_ = json.NewDecoder(request.Body).Decode(&payload)
	}

	repoRoot := payload.RepoRoot
	if repoRoot == "" {
		repoRoot = s.repoRoot
	}
	if repoRoot == "" {
		if cwd, err := os.Getwd(); err == nil {
			repoRoot = filepath.Clean(filepath.Join(cwd, "../.."))
		}
	}

	documents, err := s.service.IndexSpecs(request.Context(), repoRoot)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{"count": len(documents), "documents": documents})
}

func (s *Server) handleSearch(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeMethodNotAllowed(writer)
		return
	}
	var payload struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	bundle, err := s.service.Search(request.Context(), payload.Query, payload.Limit)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	writeJSON(writer, http.StatusOK, bundle)
}

func (s *Server) handleDraft(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeMethodNotAllowed(writer)
		return
	}
	var payload struct {
		Query  string `json:"query"`
		Title  string `json:"title"`
		Format string `json:"format"`
		Limit  int    `json:"limit"`
	}
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	draft, err := s.service.DraftSpec(request.Context(), payload.Query, payload.Title, payload.Format, payload.Limit)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	writeJSON(writer, http.StatusOK, draft)
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeError(writer http.ResponseWriter, status int, err error) {
	writeJSON(writer, status, map[string]string{"error": err.Error()})
}

func writeMethodNotAllowed(writer http.ResponseWriter) {
	writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
