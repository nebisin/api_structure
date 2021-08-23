# build stage
FROM golang:1.17-alpine AS builder
RUN apk --update add ca-certificates
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api/main.go

RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY wait-for.sh .
COPY start.sh .
COPY ./migrations ./migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]