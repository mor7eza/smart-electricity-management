FROM golang:1.24.6-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o transmitter .

FROM alpine:3.9.6

WORKDIR /root/

COPY --from=builder /app/transmitter .

CMD ["sh", "-c", "./transmitter --count=$COUNT --interval=$INTERVAL"]
