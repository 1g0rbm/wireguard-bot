# Первый этап - сборка
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/bot.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

CMD ["./main"]
