# Сборка внутри образа: не зависит от ОС хоста (Windows/Linux/macOS).
FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/weighing ./cmd/api

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /out/weighing /app/weighing

EXPOSE 8089

CMD ["/app/weighing"]
