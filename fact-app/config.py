import os

class Config:
    HTTP_PORT = int(os.environ.get("HTTP_PORT", "8001"))
    TENANT_DB_CONN = os.environ.get("TENANT_DB_CONN", "postgres://postgres:postgres@localhost:5432/fact_app_tenant_registry?sslmode=disable")
    # Base connection for creating tenant DBs (same host, DB name = fact-app_<tenant_id>)
    DB_HOST = os.environ.get("DB_HOST", "localhost")
    DB_PORT = os.environ.get("DB_PORT", "5432")
    DB_USER = os.environ.get("DB_USER", "postgres")
    DB_PASSWORD = os.environ.get("DB_PASSWORD", "postgres")
    KAFKA_BOOTSTRAP = os.environ.get("KAFKA_BOOTSTRAP", "localhost:9092")
    KAFKA_TOPIC = os.environ.get("KAFKA_TOPIC", "tenant.signups")
    APP_NAME = "fact-app"
