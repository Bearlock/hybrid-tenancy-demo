# todo-app

Elixir app: CRUD for TODOs per tenant.

- **TenantDB**: Registry DB (`todo_app_tenant_registry`) stores `{id, host}` per tenant. Tenant DBs are named `todo-app_<tenant_id>`.
- **Kafka**: Optional. Run `mix run -e "TodoApp.KafkaConsumer.run()"` in a separate process to consume `tenant.signups` and create tenant DBs when `todo-app` is in `services`.
- **API**: Expects `X-Tenant-ID`. CRUD on `/todos`.

## Run

```bash
# Create registry DB first (e.g. docker/init-dbs.sh creates todo_app_tenant_registry).
mix deps.get
export DB_HOST=localhost
export TENANT_REGISTRY_DB=todo_app_tenant_registry
mix run --no-halt
```

## API (with X-Tenant-ID)

- `GET /todos` — list todos
- `POST /todos` — body `{"title": "..."}`
- `GET /todos/:id`
- `PUT /todos/:id` — body `{"title": "...", "completed": true}`
- `DELETE /todos/:id`

## Kafka consumer (optional)

In another terminal, after the app is running and Kafka is up:

```bash
mix run -e "TodoApp.KafkaConsumer.run()"
```

This creates tenant DBs when new tenants sign up for the todo-app service.
