package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hybrid-tenancy-demo/org-app/internal/db"
)

type OrgHandler struct {
	registry *db.Registry
	cfg      *struct {
		DBHost, DBPort, DBUser, DBPassword string
	}
}

func NewOrgHandler(registry *db.Registry, dbHost, dbPort, dbUser, dbPassword string) *OrgHandler {
	return &OrgHandler{
		registry: registry,
		cfg: &struct {
			DBHost, DBPort, DBUser, DBPassword string
		}{dbHost, dbPort, dbUser, dbPassword},
	}
}

func (h *OrgHandler) tenantDB(r *http.Request) (*sql.DB, error) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		return nil, nil
	}
	host, err := h.registry.Host(tenantID)
	if err != nil {
		return nil, err
	}
	return db.OpenTenantDB(host, h.cfg.DBPort, h.cfg.DBUser, h.cfg.DBPassword, tenantID)
}

type orgUnit struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id,omitempty"`
}

func (h *OrgHandler) List(w http.ResponseWriter, r *http.Request) {
	conn, err := h.tenantDB(r)
	if conn == nil {
		http.Error(w, "missing X-Tenant-ID", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}
	defer conn.Close()

	rows, err := conn.QueryContext(r.Context(), "SELECT id, name, parent_id FROM org_units ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var units []orgUnit
	for rows.Next() {
		var u orgUnit
		var pID sql.NullInt64
		if err := rows.Scan(&u.ID, &u.Name, &pID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if pID.Valid {
			id := int(pID.Int64)
			u.ParentID = &id
		}
		units = append(units, u)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(units)
}

func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	conn, err := h.tenantDB(r)
	if conn == nil {
		http.Error(w, "missing X-Tenant-ID", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}
	defer conn.Close()

	var body struct {
		Name     string `json:"name"`
		ParentID *int   `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if body.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	var id int
	if body.ParentID != nil {
		err = conn.QueryRowContext(r.Context(), "INSERT INTO org_units (name, parent_id) VALUES ($1, $2) RETURNING id", body.Name, *body.ParentID).Scan(&id)
	} else {
		err = conn.QueryRowContext(r.Context(), "INSERT INTO org_units (name) VALUES ($1) RETURNING id", body.Name).Scan(&id)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(orgUnit{ID: id, Name: body.Name, ParentID: body.ParentID})
}

func (h *OrgHandler) Get(w http.ResponseWriter, r *http.Request) {
	conn, err := h.tenantDB(r)
	if conn == nil {
		http.Error(w, "missing X-Tenant-ID", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}
	defer conn.Close()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var u orgUnit
	var pID sql.NullInt64
	err = conn.QueryRowContext(r.Context(), "SELECT id, name, parent_id FROM org_units WHERE id = $1", id).Scan(&u.ID, &u.Name, &pID)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pID.Valid {
		pid := int(pID.Int64)
		u.ParentID = &pid
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)
}

func (h *OrgHandler) Update(w http.ResponseWriter, r *http.Request) {
	conn, err := h.tenantDB(r)
	if conn == nil {
		http.Error(w, "missing X-Tenant-ID", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}
	defer conn.Close()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body struct {
		Name     string `json:"name"`
		ParentID *int   `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if body.ParentID != nil {
		_, err = conn.ExecContext(r.Context(), "UPDATE org_units SET name = $1, parent_id = $2 WHERE id = $3", body.Name, *body.ParentID, id)
	} else {
		_, err = conn.ExecContext(r.Context(), "UPDATE org_units SET name = $1, parent_id = NULL WHERE id = $2", body.Name, id)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u := orgUnit{ID: id, Name: body.Name, ParentID: body.ParentID}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)
}

func (h *OrgHandler) Delete(w http.ResponseWriter, r *http.Request) {
	conn, err := h.tenantDB(r)
	if conn == nil {
		http.Error(w, "missing X-Tenant-ID", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}
	defer conn.Close()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	_, err = conn.ExecContext(r.Context(), "DELETE FROM org_units WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
