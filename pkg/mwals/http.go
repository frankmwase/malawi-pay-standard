package mwals

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler provides HTTP endpoints for alias resolution.
type Handler struct {
	resolver Resolver
}

func NewHandler(r Resolver) *Handler {
	return &Handler{resolver: r}
}

// ServeHTTP handles the /resolve/@alias request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple path parsing: /resolve/@alias
	path := strings.TrimPrefix(r.URL.Path, "/resolve/")
	if path == "" {
		http.Error(w, "Missing alias", http.StatusBadRequest)
		return
	}

	resp, err := h.resolver.Resolve(r.Context(), path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
