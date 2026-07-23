# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOGC=20 go build -o server .

# Runtime stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8081

CMD ["./server"]
