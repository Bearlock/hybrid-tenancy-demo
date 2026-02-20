defmodule TodoApp.Router do
  use Plug.Router
  plug Plug.Logger
  plug :match
  plug Plug.Parsers, parsers: [:json], json_decoder: Jason
  plug :dispatch

  def get_tenant_id(conn) do
    Plug.Conn.get_req_header(conn, "x-tenant-id") |> List.first()
  end

  def tenant_conn(conn) do
    case get_tenant_id(conn) do
      nil -> nil
      tid -> TodoApp.TenantDB.get_tenant_conn(tid)
    end
  end

  get "/todos" do
    case tenant_conn(conn) do
      nil -> send_resp(conn, 401, Jason.encode!(%{error: "missing X-Tenant-ID"}))
      {:ok, db} ->
        try do
          {:ok, %{rows: rows}} = Postgrex.query(db, "SELECT id, title, completed, created_at FROM todos ORDER BY id", [])
          body = for [id, title, completed, created_at] <- rows, do: %{id: id, title: title, completed: completed, created_at: to_string(created_at)}
          conn |> put_resp_content_type("application/json") |> send_resp(200, Jason.encode!(body))
        after
          GenServer.stop(db)
        end
    end
  end

  post "/todos" do
    case tenant_conn(conn) do
      nil -> send_resp(conn, 401, Jason.encode!(%{error: "missing X-Tenant-ID"}))
      {:ok, db} ->
        title = get_in(conn.body_params, ["title"]) || ""
        if title == "" do
          send_resp(conn, 400, Jason.encode!(%{error: "title required"}))
        else
          try do
            {:ok, %{rows: [[id, title, completed, created_at] | _]}} =
              Postgrex.query(db, "INSERT INTO todos (title) VALUES ($1) RETURNING id, title, completed, created_at", [title])
            conn |> put_resp_content_type("application/json") |> put_resp_header("content-type", "application/json") |> send_resp(201, Jason.encode!(%{id: id, title: title, completed: completed, created_at: to_string(created_at)}))
          after
            GenServer.stop(db)
          end
        end
    end
  end

  get "/todos/:id" do
    id = String.to_integer(conn.path_params["id"])
    case tenant_conn(conn) do
      nil -> send_resp(conn, 401, Jason.encode!(%{error: "missing X-Tenant-ID"}))
      {:ok, db} ->
        try do
          case Postgrex.query(db, "SELECT id, title, completed, created_at FROM todos WHERE id = $1", [id]) do
            {:ok, %{rows: []}} -> send_resp(conn, 404, Jason.encode!(%{error: "not found"}))
            {:ok, %{rows: [[id, title, completed, created_at] | _]}} ->
              conn |> put_resp_content_type("application/json") |> send_resp(200, Jason.encode!(%{id: id, title: title, completed: completed, created_at: to_string(created_at)}))
          end
        after
          GenServer.stop(db)
        end
    end
  end

  put "/todos/:id" do
    id = String.to_integer(conn.path_params["id"])
    case tenant_conn(conn) do
      nil -> send_resp(conn, 401, Jason.encode!(%{error: "missing X-Tenant-ID"}))
      {:ok, db} ->
        title = get_in(conn.body_params, ["title"])
        completed = get_in(conn.body_params, ["completed"])
        try do
          if title != nil do
            Postgrex.query(db, "UPDATE todos SET title = $1 WHERE id = $2", [title, id])
          end
          if completed != nil do
            Postgrex.query(db, "UPDATE todos SET completed = $1 WHERE id = $2", [completed, id])
          end
          case Postgrex.query(db, "SELECT id, title, completed, created_at FROM todos WHERE id = $1", [id]) do
            {:ok, %{rows: []}} -> send_resp(conn, 404, Jason.encode!(%{error: "not found"}))
            {:ok, %{rows: [row]}} ->
              [id, title, completed, created_at] = row
              conn |> put_resp_content_type("application/json") |> send_resp(200, Jason.encode!(%{id: id, title: title, completed: completed, created_at: to_string(created_at)}))
          end
        after
          GenServer.stop(db)
        end
    end
  end

  delete "/todos/:id" do
    id = String.to_integer(conn.path_params["id"])
    case tenant_conn(conn) do
      nil -> send_resp(conn, 401, Jason.encode!(%{error: "missing X-Tenant-ID"}))
      {:ok, db} ->
        try do
          Postgrex.query(db, "DELETE FROM todos WHERE id = $1", [id])
          send_resp(conn, 204, "")
        after
          GenServer.stop(db)
        end
    end
  end

  match _ do
    send_resp(conn, 404, "Not found")
  end
end
