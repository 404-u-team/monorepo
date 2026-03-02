# TODO: переделать для этого проекта

# --- Build Stage ---
FROM golang:1.25.5-alpine3.23 AS builder

ENV CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	GOMODCACHE=/go/pkg/mod

WORKDIR /app/services

COPY services/proto ./proto

WORKDIR /app

COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/

WORKDIR /app/services/auth-service
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY services/auth-service .

RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go build -trimpath -ldflags="-s -w" -o server ./cmd/api/main.go


# --- Run Stage ---
FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/services/auth-service/server .

COPY --from=builder /app/services/auth-service/internal/db/migrations ./migrations

# создаем обычного пользователя (не root) и используем его
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["./server"]