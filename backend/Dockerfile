# ─────────────────────────────────────────────────────────
#  BidFlow Backend – Dockerfile
#  Produces: Go API server binary
# ─────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

# ── Runtime ──────────────────────────────────────────────
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/server ./server
EXPOSE 8080
CMD ["./server"]
