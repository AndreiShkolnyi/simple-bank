#Build stage
FROM golang:1.24.5-alpine3.22 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

#Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
# COPY app.env .
COPY start.sh .
COPY db/migration ./migration

EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT ["sh", "/app/start.sh"]