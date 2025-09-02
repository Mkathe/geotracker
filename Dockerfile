# syntax=docker/dockerfile:1.4
FROM golang:tip-bullseye AS builder
LABEL maintainer="Magzhan Yedilbayev"

WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1
ENV MIGRATE_VERSION v4.17.0

RUN apt-get update && apt-get install -y wget ca-certificates && \
    wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz && \
    tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz && \
    rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/$MIGRATE_VERSION/migrate.linux-amd64.tar.gz && \
    tar -xvzf migrate.tar.gz -C /usr/local/bin && \
    rm migrate.tar.gz

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

COPY ./pkg/db/migrations /app/migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates postgresql-client && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /usr/local/bin/dockerize /usr/local/bin/dockerize
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["/bin/sh", "-c", "\
  dockerize \
    -wait tcp://postgres_library:5432 \
    -wait tcp://keycloak_library:8080 \
    -wait-retry-interval 5s \
    -timeout 120s && \
  until pg_isready -h postgres_library -p 5432 -U ${DB_USER}; do \
    echo 'Waiting for Postgres to be ready...'; \
    sleep 2; \
  done && \
  migrate -path /app/migrations -database ${CONN_DB_POSTGRES} up && \
  exec /app/main"]

#-wait tcp://kafka_library:9092