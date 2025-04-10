FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-task-tracker .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/go-task-tracker .
COPY --from=builder /app/adapter/repo/postgres/migrations ./adapter/repo/postgres/migrations

RUN chmod +x go-task-tracker

EXPOSE ${API_PORT}

ENV API_PORT=8080
ENV DB_CONNECTION_URL="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

CMD ["./go-task-tracker"]
