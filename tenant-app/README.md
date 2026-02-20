# tenant-app

Go service that owns **TenantMetaDB** and is responsible for:

- **Tenants**: Records tenant `id` and `name` in the `tenants` table.
- **Services**: Records which services a tenant signed up for (`fact-app`, `org-app`, `todo-app`) in `tenant_services`.
- **API tokens**: Mints a long-lived JWT for Bearer auth with the API gateway; stores token hash in `api_tokens`.
- **Events**: On successful signup, publishes a `tenant.signups` Kafka event so downstream apps (fact-app, org-app, todo-app) can create tenant databases.

All endpoints are under the `/tenants` path.

## Run

```bash
# Ensure TenantMetaDB and Kafka are up (see repo docker-compose / compose).
export META_DB_CONN="postgres://postgres:postgres@localhost:5432/tenant_meta?sslmode=disable"
export KAFKA_BROKERS="localhost:9092"
go run ./cmd/tenant-app
```

## API

All under `/tenants`:

- `POST /tenants` — Create (sign up) a tenant. Body: `{"name": "Acme", "services": ["fact-app", "org-app"]}`. Returns `{"tenant_id": "...", "token": "..."}`.
- `GET /tenants/:id` — Fetch tenant info: `{"id": "...", "name": "...", "services": ["fact-app", ...]}`. 404 if not found.
