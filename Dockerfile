FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/ledger .

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    libc6 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/ledger .

COPY .env /app/ledger/.env

RUN mkdir -p /app/logs && touch /app/logs/app.log && chmod -R 777 /app/logs

EXPOSE 8080

CMD ["./ledger"]