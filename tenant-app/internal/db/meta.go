package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const schema = `
CREATE TABLE IF NOT EXISTS tenants (
	id   TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tenant_services (
	tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
	service   TEXT NOT NULL CHECK (service IN ('fact-app', 'org-app', 'todo-app')),
	PRIMARY KEY (tenant_id, service)
);

CREATE TABLE IF NOT EXISTS api_tokens (
	id        TEXT PRIMARY KEY,
	tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_api_tokens_tenant ON api_tokens(tenant_id);
`

func OpenMetaDB(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("schema: %w", err)
	}
	return db, nil
}
