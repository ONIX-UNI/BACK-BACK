FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/worker ./cmd/workers/mailer

FROM alpine:3.21 AS runtime
RUN set -eux; \
    attempts=0; \
    until [ "$attempts" -ge 5 ]; do \
      apk add --no-cache ca-certificates tzdata && break; \
      attempts=$((attempts + 1)); \
      echo "apk add failed (attempt ${attempts}/5), retrying..." >&2; \
      sleep $((attempts * 2)); \
    done; \
    [ "$attempts" -lt 5 ]
WORKDIR /app

FROM runtime AS api
COPY --from=builder /out/api /usr/local/bin/api
ENTRYPOINT ["/usr/local/bin/api"]

FROM runtime AS worker
COPY --from=builder /out/worker /usr/local/bin/worker
ENTRYPOINT ["/usr/local/bin/worker"]
