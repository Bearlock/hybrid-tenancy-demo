# Event contract: tenant signups

tenant-app publishes to Kafka when a tenant is successfully created. Downstream apps (fact-app, org-app, todo-app) consume this topic and create tenant databases when their service is in the signup payload.

## Topic

- **Name**: `tenant.signups` (configurable via `KAFKA_TOPIC`)

## Payload (JSON)

```json
{
  "tenant_id": "uuid-from-tenant-app",
  "name": "Tenant display name",
  "services": ["fact-app", "org-app", "todo-app"]
}
```

- `tenant_id`: Unique ID for the tenant (from tenant-app).
- `name`: Tenant name (for logging/audit).
- `services`: Subset of `["fact-app", "org-app", "todo-app"]` that the tenant signed up for.

## Consumer behavior

- Each app (fact-app, org-app, todo-app) subscribes to `tenant.signups`.
- If `services` contains that app’s name (`fact-app`, `org-app`, or `todo-app`), the app:
  1. Creates a logical database named `<app-name>_<tenant_id>` on the shared cluster.
  2. Applies the app’s schema in that database.
  3. Registers the tenant in its own TenantDB (registry) with `{ "id": "<tenant_id>", "host": "<db_host>" }`.

## TenantDB registry record

Each app maintains a registry of tenants it serves. A record is a logical mapping:

```json
{
  "id": "123",
  "host": "somehost"
}
```

- `id`: Tenant ID (same as `tenant_id` in the event).
- `host`: Database host for that tenant’s logical DB (same cluster in this demo; in a multi-region setup it could be a hostname or pooler address).
