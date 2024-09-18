
FROM golang:1.20 AS builder  


WORKDIR /app  


COPY go.mod go.sum ./  


RUN go mod download  


COPY . .  


RUN CGO_ENABLED=0 GOOS=linux go build -o botus ./cmd/main.go  


FROM alpine:latest  


RUN apk --no-cache add ca-certificates  


COPY --from=builder /app/botus /botus  

ENTRYPOINT ["/botus"]  