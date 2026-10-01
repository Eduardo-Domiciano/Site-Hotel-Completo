FROM golang:1.26-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata \
  && adduser -D -u 10001 appuser \
  && mkdir -p /app/uploads \
  && chown -R appuser:appuser /app

COPY --from=builder /out/api /app/api
COPY migrations /app/migrations
# SQL do host chega como 640; o appuser precisa ler as migrations.
RUN chmod -R a+rX /app/api /app/migrations

USER appuser
EXPOSE 8080

CMD ["/app/api"]
