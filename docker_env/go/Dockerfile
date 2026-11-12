# Runtime-образ. Бинарник собирается на хосте: go build -o weighing ./cmd/api
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

EXPOSE 8089

CMD ["./weighing"]
