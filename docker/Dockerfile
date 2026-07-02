FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o server ./cmd/server

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o gateway ./cmd/gateway

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/gateway .

COPY --from=builder /app/.env* ./

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 50051 50052

HEALTHCHECK --interval=60s --timeout=5s --start-period=10s --retries=3 \
    CMD nc -z localhost 50051 || exit 1

CMD ["./server"]