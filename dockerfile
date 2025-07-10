#Build stage
FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

# Final stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./
COPY app.env .
COPY start.sh .
COPY wait-for .
COPY db/migration/ ./migration
RUN chmod +x ./start.sh
RUN chmod +x ./wait-for

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh"] 