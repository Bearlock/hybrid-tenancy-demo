package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/auth"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/events"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/repo"
)

var validServices = map[string]bool{
	"fact-app": true,
	"org-app":  true,
	"todo-app": true,
}

type SignupRequest struct {
	Name     string   `json:"name"`
	Services []string `json:"services"`
}

type SignupResponse struct {
	TenantID string `json:"tenant_id"`
	Token   string `json:"token"`
}

type SignupHandler struct {
	repo     *repo.Repo
	producer *events.Producer
	key      string
}

func NewSignupHandler(repo *repo.Repo, producer *events.Producer, signingKey string) *SignupHandler {
	return &SignupHandler{repo: repo, producer: producer, key: signingKey}
}

func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	services := make([]string, 0, len(req.Services))
	for _, s := range req.Services {
		if validServices[s] {
			services = append(services, s)
		}
	}
	if len(services) == 0 {
		http.Error(w, "at least one valid service required (fact-app, org-app, todo-app)", http.StatusBadRequest)
		return
	}

	tenantID := uuid.New().String()
	if err := h.repo.CreateTenant(r.Context(), tenantID, req.Name, services); err != nil {
		http.Error(w, "failed to create tenant", http.StatusInternalServerError)
		return
	}

	rawToken, err := auth.MintToken(tenantID, services, h.key)
	if err != nil {
		http.Error(w, "failed to mint token", http.StatusInternalServerError)
		return
	}
	tokenHash := auth.HashToken(rawToken)
	tokenID := uuid.New().String()
	if err := h.repo.StoreToken(r.Context(), tokenID, tenantID, tokenHash); err != nil {
		http.Error(w, "failed to store token", http.StatusInternalServerError)
		return
	}

	if err := h.producer.PublishTenantSignup(r.Context(), events.TenantSignupEvent{
		TenantID: tenantID,
		Name:     req.Name,
		Services: services,
	}); err != nil {
		// Log but don't fail signup; downstream can be eventually consistent
		_ = err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(SignupResponse{TenantID: tenantID, Token: rawToken})
}
