FROM golang:1.26.4-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/synclink ./cmd/server

FROM alpine:3.24

RUN apk add --no-cache ca-certificates \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /bin/synclink /app/synclink
RUN mkdir -p /var/data && chown -R appuser:appgroup /app /var/data

USER appuser
EXPOSE 8080
ENV ADDR=:8080
ENV SYNCLINK_DB=/var/data/synclink.db

ENTRYPOINT ["/app/synclink"]
