FROM golang:1.24-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main ./cmd/api-gateway/main.go

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
