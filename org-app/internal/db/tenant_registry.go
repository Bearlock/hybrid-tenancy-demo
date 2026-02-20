package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const registrySchema = `
CREATE TABLE IF NOT EXISTS tenants (
	id   TEXT PRIMARY KEY,
	host TEXT NOT NULL
);
`

func OpenTenantRegistry(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Printf("Unable to open connection to DB: %s", conn)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Printf("Unable to ping DB: %s", conn)
		return nil, err
	}
	if _, err := db.Exec(registrySchema); err != nil {
		log.Printf("Unable to create registry table")
		return nil, fmt.Errorf("registry schema: %w", err)
	}
	return db, nil
}

func RegisterTenant(db *sql.DB, tenantID, host string) error {
	_, err := db.Exec(
		`INSERT INTO tenants (id, host) VALUES ($1, $2)`,
		tenantID, host)

	if err != nil {
		log.Printf("failed to insert tenant into tenant table: %s", tenantID)
	}
	log.Printf("inserted tenant into tenant table: %s", tenantID)
	return err
}

func GetTenantHost(db *sql.DB, tenantID string) (host string, err error) {
	err = db.QueryRow("SELECT host FROM tenants WHERE id = $1", tenantID).Scan(&host)
	return host, err
}

type Registry struct {
	db *sql.DB
}

func NewRegistry(db *sql.DB) *Registry {
	return &Registry{db: db}
}

func (r *Registry) Register(tenantID, host string) error {
	return RegisterTenant(r.db, tenantID, host)
}

func (r *Registry) Host(tenantID string) (string, error) {
	return GetTenantHost(r.db, tenantID)
}
