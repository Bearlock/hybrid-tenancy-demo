defmodule TodoApp.KafkaSubscriber do
  use Broadway

  require Logger

  @app_name "todo-app"

  def start_link(_opts) do
    Broadway.start_link(__MODULE__,
      name: __MODULE__,
      producer: [
        module:
          {
            BroadwayKafka.Producer,
            [
              hosts: [kafka: 29092],
              group_id: @app_name,
              topics: ["tenant.signups"]
            ]
          },
        concurrency: 1
      ],
      processors: [
        default: [
          concurrency: 10
        ]
      ]
    )
  end

  @impl true
  def handle_message(_, message, _) do
    IO.inspect(message, label: "message")

    case Jason.decode(message.data) do
      {:ok, %{"tenant_id" => tid, "services" => services}} when is_list(services) ->
        if @app_name in services do
          case TodoApp.TenantDB.create_tenant_database(tid) do
            :ok -> Logger.info("Created tenant DB for #{tid}")
            err -> Logger.warning("Failed to create tenant DB for #{tid}: #{inspect(err)}")
          end
        end
      _ -> :ok
    end
    message
  end
end
