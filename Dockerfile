
FROM golang:1.18 AS builder


WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .

# Собираем приложение.
RUN go build -o botus ./cmd/main.go


FROM alpine:latest


WORKDIR /root/


COPY --from=builder /app/botus .


CMD ["./botus"]