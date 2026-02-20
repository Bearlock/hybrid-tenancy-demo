defmodule TodoApp.TenantDB do
  @moduledoc "Create tenant DB and open connections. DB name = todo-app_<tenant_id>"
  require Logger

  def db_name(tenant_id), do: "todo-app_#{tenant_id}"

  def create_tenant_database(tenant_id) do
    cfg = Application.get_env(:todo_app, :tenant_registry)
    host = to_charlist(cfg[:hostname])
    port = cfg[:port]
    user = to_charlist(cfg[:username])
    pass = to_charlist(cfg[:password] || "")
    db = db_name(tenant_id)

    opts = [hostname: host, port: port, username: user, password: pass, database: "postgres"]
    {:ok, conn} = Postgrex.start_link(opts)
    try do
      Postgrex.query!(conn, "CREATE DATABASE \"#{db}\"", [])
    rescue
      e in Postgrex.Error -> if e.postgres.code != :duplicate_database, do: reraise e, __STACKTRACE__
    end
    GenServer.stop(conn)

    schema_opts = [hostname: host, port: port, username: user, password: pass, database: db]
    {:ok, schema_conn} = Postgrex.start_link(schema_opts)
    Postgrex.query!(schema_conn, """
    CREATE TABLE IF NOT EXISTS todos (
      id SERIAL PRIMARY KEY,
      title TEXT NOT NULL,
      completed BOOLEAN DEFAULT FALSE,
      created_at TIMESTAMPTZ DEFAULT NOW()
    )
    """, [])
    GenServer.stop(schema_conn)

    TodoApp.TenantRegistry.register(tenant_id, to_string(cfg[:hostname]))
    :ok
  end

  def get_tenant_conn(tenant_id) do
    case TodoApp.TenantRegistry.get_host(tenant_id) do
      {:ok, host} -> open_conn(tenant_id, host)
      :error -> nil
    end
  end

  defp open_conn(tenant_id, host) do
    cfg = Application.get_env(:todo_app, :tenant_registry)
    opts = [
      hostname: to_charlist(host),
      port: cfg[:port],
      username: to_charlist(cfg[:username]),
      password: to_charlist(cfg[:password] || ""),
      database: db_name(tenant_id)
    ]
    case Postgrex.start_link(opts) do
      {:ok, pid} -> {:ok, pid}
      e -> e
    end
  end
end
