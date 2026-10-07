package mwals

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Handler provides HTTP endpoints for alias resolution.
type Handler struct {
	resolver          *Service
	registrationToken string
}

func NewHandler(s *Service) *Handler {
	return &Handler{resolver: s}
}

// NewAuthenticatedHandler enables registration only with the supplied bearer token.
func NewAuthenticatedHandler(s *Service, token string) *Handler {
	return &Handler{resolver: s, registrationToken: token}
}

// ServeHTTP handles the /resolve/@alias request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple path parsing: /resolve/@alias
	path := strings.TrimPrefix(r.URL.Path, "/resolve/")
	if !strings.HasPrefix(r.URL.Path, "/resolve/") || path == "" || strings.Contains(path, "/") {
		http.Error(w, "Invalid alias path", http.StatusBadRequest)
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
	if h.registrationToken == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+h.registrationToken)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var req RegistrationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if decoder.Decode(new(any)) != io.EOF {
		http.Error(w, "Unexpected trailing data", http.StatusBadRequest)
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
		// Never disclose internal persistence details in an HTTP response.
		http.Error(w, "Registration rejected", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "registered", "alias": req.Alias})
}
