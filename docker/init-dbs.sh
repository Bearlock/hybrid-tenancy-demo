#!/usr/bin/env bash
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE tenant_meta;
  CREATE DATABASE fact_app_tenant_registry;
  CREATE DATABASE org_app_tenant_registry;
  CREATE DATABASE todo_app_tenant_registry;
EOSQL
