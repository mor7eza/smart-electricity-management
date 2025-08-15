FROM golang:1.24.6-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o billing-service ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/billing-service .

CMD ["./billing-service"]
