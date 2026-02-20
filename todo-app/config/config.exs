import Config

config :todo_app, :http_port, String.to_integer(System.get_env("HTTP_PORT") || "8003")

config :todo_app, :tenant_registry,
  hostname: System.get_env("DB_HOST") || "localhost",
  port: String.to_integer(System.get_env("DB_PORT") || "5432"),
  username: System.get_env("DB_USER") || "postgres",
  password: System.get_env("DB_PASSWORD") || "postgres",
  database: System.get_env("TENANT_REGISTRY_DB") || "todo_app_tenant_registry"

config :todo_app, :kafka,
  brokers: [{"kafka", 29092}],
  topic: "tenant.signups",
  group_id: "todo-app"

config :brod,
  clients: [
    kafka_clients: [
      endpoints: [kafka: 29092]
    ]
  ]

config :todo_app, :app_name, "todo-app"
