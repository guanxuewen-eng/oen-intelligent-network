FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn,direct go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /oen-server ./cmd/server

# Final stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /oen-server /app/oen-server
COPY --from=builder /build/config.yaml /app/config.yaml

EXPOSE 8080

ENTRYPOINT ["/app/oen-server"]
CMD ["-config", "/app/config.yaml"]
