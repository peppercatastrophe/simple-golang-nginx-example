# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOGC=20 go build -o server .

# Runtime stage
FROM alpine:latest AS runtime

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8081

CMD ["./server"]
