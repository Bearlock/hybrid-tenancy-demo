# org-app

Go service: lightweight org chart per tenant.

- **TenantDB**: Registry DB (`org_app_tenant_registry`) stores `{id, host}` per tenant. Tenant DBs are named `org-app_<tenant_id>`.
- **Kafka**: Consumes `tenant.signups`; when `org-app` is in `services`, creates the tenant DB and registers it.
- **API**: Expects `X-Tenant-ID`. CRUD for org units (name, optional parent_id).

## Run

```bash
export TENANT_DB_CONN="postgres://postgres:postgres@localhost:5432/org_app_tenant_registry?sslmode=disable"
export KAFKA_BROKERS="localhost:9092"
go run ./cmd/org-app
```

## API (with X-Tenant-ID)

- `GET /org` — list org units
- `POST /org` — body `{"name": "...", "parent_id": 1}` (parent_id optional)
- `GET /org/{id}`
- `PUT /org/{id}`
- `DELETE /org/{id}`
