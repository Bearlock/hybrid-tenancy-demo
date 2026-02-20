# Hybrid tenancy demo

Lightweight proof-of-concept for **hybrid tenancy**: each tenant has its own **logical database** on a **shared database cluster**. Each app has a **TenantDB** registry (`id` → `host`) and creates DBs named `<app>_<tenant_id>`.

## Architecture

See **[docs/architecture.md](docs/architecture.md)** for a diagram and flow description.

| App | Role | Stack |
|-----|------|--------|
| **tenant-app** | Tenant signup & info, TenantMetaDB, API token minting, Kafka producer | Go |
| **api-gateway-app** | API gateway, Bearer JWT auth, routes to internal apps | Go |
| **fact-app** | CRUD random facts per tenant | Python (Flask) |
| **org-app** | CRUD org chart per tenant | Go |
| **todo-app** | CRUD TODOs per tenant | Elixir |

- **TenantMetaDB** (tenant-app): `tenants`, `tenant_services`, `api_tokens`.
- **TenantDB** (per app): Registry of `{ "id": "<tenant_id>", "host": "<host>" }`; tenant DBs named `<application>_<tenant_id>`.
- **Events**: tenant-app publishes to Kafka topic `tenant.signups`. fact-app, org-app, and todo-app consume it and create tenant DBs when their service is in the signup.

## Quick start

1. **Start infrastructure**

   ```bash
   docker-compose up -d
   # Wait for Postgres and Kafka to be ready
   ```
2. **Flow**

   - `POST /tenants` (tenant-app) with `{"name": "Acme", "services": ["fact-app", "todo-app"]}` → returns `tenant_id` and `token`.
   - Use `Authorization: Bearer <token>` with the gateway; gateway forwards to fact-app, org-app, and todo-app with `X-Tenant-ID`.
   - After consumers run (or tenant DBs created otherwise), CRUD on `/facts`, `/todos`, `/org` is per-tenant.

## Repo layout

- `tenant-app/` — TenantMetaDB, signup & tenant info API (under `/tenants`), token minting, Kafka producer
- `api-gateway-app/` — Gateway, JWT auth, proxy to fact/org/todo
- `fact-app/` — Flask facts API, TenantDB, Kafka consumer
- `org-app/` — Go org API, TenantDB, Kafka consumer
- `todo-app/` — Elixir todos API, TenantDB, optional Kafka consumer
- `docker-compose.yml` — Postgres + Kafka
- `docs/events.md` — Event contract for `tenant.signups`

Each app’s README has run and API details.
