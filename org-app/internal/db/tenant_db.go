package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func TenantDBName(tenantID string) string {
	return "org-app_" + tenantID
}

const orgSchema = `
CREATE TABLE IF NOT EXISTS org_units (
	id       SERIAL PRIMARY KEY,
	name     TEXT NOT NULL,
	parent_id INT REFERENCES org_units(id),
	created_at TIMESTAMPTZ DEFAULT NOW()
);
`

func CreateTenantDatabase(host, port, user, password, tenantID string) error {
	dbName := TenantDBName(tenantID)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec("CREATE DATABASE \"" + dbName + "\""); err != nil {
		if !isDuplicateDB(err) {
			return err
		}
	}
	appConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	appDB, err := sql.Open("postgres", appConnStr)
	if err != nil {
		return err
	}
	defer appDB.Close()
	_, err = appDB.Exec(orgSchema)
	return err
}

func isDuplicateDB(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "already exists") || strings.Contains(s, "duplicate")
}

// OpenTenantDB opens a connection to the logical DB for the given tenant.
func OpenTenantDB(host, port, user, password, tenantID string) (*sql.DB, error) {
	dbName := TenantDBName(tenantID)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	return sql.Open("postgres", connStr)
}
