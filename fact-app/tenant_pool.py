"""
Resolve tenant ID to a DB connection for the tenant's logical database.
DB name = fact-app_<tenant_id>
"""
import psycopg2
from config import Config
from tenant_db import get_tenant_host, register_tenant

def tenant_db_name(tenant_id):
    return f"fact-app_{tenant_id}"

def get_tenant_conn(tenant_id):
    row = get_tenant_host(Config.TENANT_DB_CONN, tenant_id)
    if not row:
        return None
    host = row["host"]
    return psycopg2.connect(
        host=host,
        port=Config.DB_PORT,
        user=Config.DB_USER,
        password=Config.DB_PASSWORD,
        dbname=tenant_db_name(tenant_id),
    )

def create_tenant_database(tenant_id, host=None):
    if host is None:
        host = Config.DB_HOST
    db_name = tenant_db_name(tenant_id)
    # Connect to default DB to create the new database
    conn = psycopg2.connect(
        host=host,
        port=Config.DB_PORT,
        user=Config.DB_USER,
        password=Config.DB_PASSWORD,
        dbname="postgres",
    )
    conn.autocommit = True
    try:
        with conn.cursor() as cur:
            cur.execute(f'CREATE DATABASE "{db_name}"')
    except psycopg2.errors.DuplicateDatabase:
        pass
    finally:
        conn.close()
    # Create schema in the new DB
    schema_conn = psycopg2.connect(
        host=host,
        port=Config.DB_PORT,
        user=Config.DB_USER,
        password=Config.DB_PASSWORD,
        dbname=db_name,
    )
    with schema_conn.cursor() as cur:
        cur.execute("""
            CREATE TABLE IF NOT EXISTS facts (
                id SERIAL PRIMARY KEY,
                content TEXT NOT NULL,
                created_at TIMESTAMPTZ DEFAULT NOW()
            )
        """)
    schema_conn.commit()
    schema_conn.close()
    register_tenant(Config.TENANT_DB_CONN, tenant_id, host)
