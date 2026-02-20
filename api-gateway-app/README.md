# api-gateway-app

Go API Gateway: single entrypoint for external clients.

- **Auth**: Validates Bearer JWT (same signing key as tenant-app). Token payload includes `tenant_id` and `services`.
- **Routing**: Forwards to internal apps by path prefix:
  - `/facts/*` → fact-app
  - `/org/*` → org-app
  - `/todos/*` → todo-app
- Injects `X-Tenant-ID` from the JWT so backend apps can resolve the tenant DB.

## Run

```bash
export JWT_SIGNING_KEY="change-me-in-production"  # same as tenant-app
export FACT_APP_URL="http://localhost:8001"
export ORG_APP_URL="http://localhost:8002"
export TODO_APP_URL="http://localhost:8003"
go run ./cmd/api-gateway
```

## Usage

```bash
curl -H "Authorization: Bearer <token-from-signup>" http://localhost:8000/facts
curl -H "Authorization: Bearer <token-from-signup>" http://localhost:8000/todos
curl -H "Authorization: Bearer <token-from-signup>" http://localhost:8000/org
```
