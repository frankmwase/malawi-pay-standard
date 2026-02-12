package mwals

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler provides HTTP endpoints for alias resolution.
type Handler struct {
	resolver *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{resolver: s}
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

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	record := &AliasRecord{
		Alias:        req.Alias,
		Status:       AliasStatusActive,
		IdentityMask: req.IdentityMask,
		Attestation:  AttestationUnverified,
		Endpoints:    req.Endpoints,
		IsPrivate:    req.IsPrivate,
	}

	if err := h.resolver.Register(r.Context(), record); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "registered", "alias": req.Alias})
}
