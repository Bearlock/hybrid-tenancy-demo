defmodule TodoApp.TenantRegistry do
  @moduledoc "TenantDB: registry of tenant_id -> host. DB name for a tenant = todo-app_<tenant_id>"
  use GenServer
  require Logger

  def start_link(opts \\ []) do
    GenServer.start_link(__MODULE__, opts, name: __MODULE__)
  end

  @impl true
  def init(_opts) do
    cfg = Application.get_env(:todo_app, :tenant_registry)
    db_opts = [
      hostname: cfg[:hostname],
      port: cfg[:port],
      username: cfg[:username],
      password: cfg[:password],
      database: cfg[:database]
    ]
    {:ok, pid} = Postgrex.start_link(db_opts)
    create_schema(pid)
    {:ok, %{conn: pid, cfg: cfg}}
  end

  defp create_schema(pid) do
    Postgrex.query!(pid, """
    CREATE TABLE IF NOT EXISTS tenants (
      id   TEXT PRIMARY KEY,
      host TEXT NOT NULL
    )
    """, [])
  end

  def register(tenant_id, host) do
    GenServer.call(__MODULE__, {:register, tenant_id, host})
  end

  def get_host(tenant_id) do
    GenServer.call(__MODULE__, {:get_host, tenant_id})
  end

  @impl true
  def handle_call({:register, tenant_id, host}, _from, state) do
    res = Postgrex.query(state.conn, "INSERT INTO tenants (id, host) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET host = EXCLUDED.host", [tenant_id, host])
    {:reply, res, state}
  end

  def handle_call({:get_host, tenant_id}, _from, state) do
    case Postgrex.query(state.conn, "SELECT host FROM tenants WHERE id = $1", [tenant_id]) do
      {:ok, %{rows: [[host]]}} -> {:reply, {:ok, host}, state}
      _ -> {:reply, :error, state}
    end
  end
end
