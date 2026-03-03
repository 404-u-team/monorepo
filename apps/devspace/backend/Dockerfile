# --- Build stage ---
FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

RUN apk add --no-cache git ca-certificates
WORKDIR /src

# copy modules files and download dependencies first (cacheable)
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the source
COPY . ./

# build the server binary (adjust build path if different)
RUN --mount=type=cache,target=/root/.cache/go-build \
	go build -trimpath -ldflags='-s -w' -o /app/server ./cmd/api


# --- Runtime stage ---
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app

# create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/server ./server

USER appuser

ENV PORT=8080
EXPOSE 8080
CMD ["./server"]