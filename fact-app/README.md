# fact-app

Python Flask app: CRUD for random facts per tenant.

- **TenantDB**: Registry DB (`fact_app_tenant_registry`) stores `{id, host}` per tenant. Tenant databases are named `fact-app_<tenant_id>`.
- **Kafka**: Consumes `tenant.signups`; when `fact-app` is in `services`, creates the tenant DB and registers it.
- **API**: All routes require `X-Tenant-ID` (set by api-gateway). CRUD on `/facts`.

## Run

```bash
# Create registry DB first (see docker-compose or create fact_app_tenant_registry).
pip install -r requirements.txt
export TENANT_DB_CONN="postgres://postgres:postgres@localhost:5432/fact_app_tenant_registry?sslmode=disable"
export DB_HOST=localhost
python app.py
# In another terminal (optional): run consumer to react to signups
python run_consumer.py
```

## API (with X-Tenant-ID)

- `GET /facts` — list facts
- `POST /facts` — body `{"content": "..."}`
- `GET /facts/<id>`
- `PUT /facts/<id>` — body `{"content": "..."}`
- `DELETE /facts/<id>`
