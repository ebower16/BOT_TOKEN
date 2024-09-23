FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o botus ./cmd

FROM alpine:latest

COPY --from=builder /app/botus .

ENTRYPOINT ["./botus"]
