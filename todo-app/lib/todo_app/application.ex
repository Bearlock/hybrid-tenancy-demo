defmodule TodoApp.Application do
  @moduledoc false
  use Application

  def start(_type, _args) do
    children = [
      TodoApp.TenantRegistry,
      TodoApp.KafkaSubscriber,
      {Plug.Cowboy, scheme: :http, plug: TodoApp.Router, options: [port: Application.get_env(:todo_app, :http_port)]}
    ]
    opts = [strategy: :one_for_one, name: TodoApp.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
