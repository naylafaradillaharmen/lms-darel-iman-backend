# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 menghasilkan static binary murni tanpa dependensi C, kompatibel 100% dengan Linux ARM64 (aarch64) dan AMD64 di Easypanel
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# Final stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env.example /app/.env.example

EXPOSE 8080

CMD ["/app/server"]
