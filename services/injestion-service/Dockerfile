FROM golang:1.24.6-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o injestion-service ./cmd

FROM alpine:3.9.6

WORKDIR /root/

COPY --from=builder /app/injestion-service .

CMD ["./injestion-service"]
