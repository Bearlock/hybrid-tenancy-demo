FROM elixir:1.19-alpine
WORKDIR /app

RUN apk add --no-cache \
  build-base \
  git \
  cmake

# Install Hex and Rebar (Elixir build tools)
RUN mix local.hex --force && mix local.rebar --force

# Install Elixir dependencies first for layer caching
COPY . /app
RUN mix deps.get && mix deps.compile

RUN mix compile

# Build the release
# RUN mix release

CMD ["mix", "run", "--no-halt"]
