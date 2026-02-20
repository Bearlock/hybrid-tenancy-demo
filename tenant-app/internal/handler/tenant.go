package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/repo"
)

type TenantHandler struct {
	repo *repo.Repo
}

func NewTenantHandler(repo *repo.Repo) *TenantHandler {
	return &TenantHandler{repo: repo}
}

func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "id")
	if tenantID == "" {
		http.Error(w, "tenant id required", http.StatusBadRequest)
		return
	}
	t, err := h.repo.GetTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "failed to fetch tenant", http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       t.ID,
		"name":     t.Name,
		"services": t.Services,
	})
}

func (h *TenantHandler) GetTenants(w http.ResponseWriter, r *http.Request) {
	t, err := h.repo.GetTenants(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch tenants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}
