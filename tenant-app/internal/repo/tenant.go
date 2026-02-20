package repo

import (
	"context"
	"database/sql"
)

type Tenant struct {
	ID   string
	Name string
}

type TenantService struct {
	TenantID string
	Service  string
}

func (r *Repo) CreateTenant(ctx context.Context, id, name string, services []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO tenants (id, name) VALUES ($1, $2)`, id, name); err != nil {
		return err
	}
	for _, svc := range services {
		if _, err := tx.ExecContext(ctx, `INSERT INTO tenant_services (tenant_id, service) VALUES ($1, $2)`, id, svc); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repo) StoreToken(ctx context.Context, id, tenantID, tokenHash string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO api_tokens (id, tenant_id, token_hash) VALUES ($1, $2, $3)
		 ON CONFLICT (tenant_id) DO UPDATE SET token_hash = EXCLUDED.token_hash, id = EXCLUDED.id`,
		id, tenantID, tokenHash)
	return err
}

// TenantWithServices is tenant row plus the list of services they signed up for.
type TenantWithServices struct {
	ID       string
	Name     string
	Services []string
}

func (r *Repo) GetTenant(ctx context.Context, tenantID string) (*TenantWithServices, error) {
	var id, name string
	if err := r.db.QueryRowContext(ctx, `SELECT id, name FROM tenants WHERE id = $1`, tenantID).Scan(&id, &name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT service FROM tenant_services WHERE tenant_id = $1 ORDER BY service`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var services []string
	for rows.Next() {
		var svc string
		if err := rows.Scan(&svc); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return &TenantWithServices{ID: id, Name: name, Services: services}, nil
}

func (r *Repo) GetTenants(ctx context.Context) ([]TenantWithServices, error) {
	tenantRows, err := r.db.QueryContext(ctx, `SELECT id, name FROM tenants ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer tenantRows.Close()

	var tenants []TenantWithServices
	for tenantRows.Next() {
		var id, name string
		if err := tenantRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		tenants = append(tenants, TenantWithServices{ID: id, Name: name, Services: nil})
	}
	if err := tenantRows.Err(); err != nil {
		return nil, err
	}

	svcRows, err := r.db.QueryContext(ctx, `SELECT tenant_id, service FROM tenant_services ORDER BY tenant_id, service`)
	if err != nil {
		return nil, err
	}
	defer svcRows.Close()

	servicesByTenant := make(map[string][]string)
	for svcRows.Next() {
		var tenantID, svc string
		if err := svcRows.Scan(&tenantID, &svc); err != nil {
			return nil, err
		}
		servicesByTenant[tenantID] = append(servicesByTenant[tenantID], svc)
	}
	if err := svcRows.Err(); err != nil {
		return nil, err
	}

	for i := range tenants {
		tenants[i].Services = servicesByTenant[tenants[i].ID]
		if tenants[i].Services == nil {
			tenants[i].Services = []string{}
		}
	}
	return tenants, nil
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}
