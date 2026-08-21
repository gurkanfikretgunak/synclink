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
RUN chown appuser:appgroup /app/synclink

USER appuser
EXPOSE 8080
ENV ADDR=:8080

ENTRYPOINT ["/app/synclink"]
