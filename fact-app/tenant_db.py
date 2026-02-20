"""
TenantDB: registry of tenant_id -> host (and optional connection info).
Each record: {"id": "123", "host": "somehost"}
"""
import psycopg2
import psycopg2.extras
from contextlib import contextmanager
import os

TENANT_REGISTRY_SCHEMA = """
CREATE TABLE IF NOT EXISTS tenants (
    id   TEXT PRIMARY KEY,
    host TEXT NOT NULL
);
"""

def get_tenant_registry_conn(conn_str):
    return psycopg2.connect(conn_str)

def init_tenant_registry(conn):
    with conn.cursor() as cur:
        cur.execute(TENANT_REGISTRY_SCHEMA)
    conn.commit()

def ensure_tenant_registry(conn_str):
    conn = get_tenant_registry_conn(conn_str)
    init_tenant_registry(conn)
    conn.close()

def register_tenant(conn_str, tenant_id, host):
    conn = get_tenant_registry_conn(conn_str)
    try:
        with conn.cursor() as cur:
            cur.execute(
                "INSERT INTO tenants (id, host) VALUES (%s, %s) ON CONFLICT (id) DO UPDATE SET host = EXCLUDED.host",
                (tenant_id, host)
            )
        conn.commit()
    finally:
        conn.close()

def get_tenant_host(conn_str, tenant_id):
    conn = get_tenant_registry_conn(conn_str)
    try:
        with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
            cur.execute("SELECT id, host FROM tenants WHERE id = %s", (tenant_id,))
            row = cur.fetchone()
            return row
    finally:
        conn.close()

def list_tenants(conn_str):
    conn = get_tenant_registry_conn(conn_str)
    try:
        with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
            cur.execute("SELECT id, host FROM tenants")
            return cur.fetchall()
    finally:
        conn.close()
