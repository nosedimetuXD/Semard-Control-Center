# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server cmd/api/main.go

# Production Minimal Stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S semard && adduser -S semard -G semard

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations /app/migrations

RUN mkdir -p /app/storage && chown -R semard:semard /app

USER semard

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/app/server"]
