FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o sso ./cmd/sso/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sso .
COPY --from=builder /app/config ./config



RUN ls -la /app/config

CMD ["./sso"]