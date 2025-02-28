FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go


FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache postgresql-client
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /app/goose
COPY .env wait-for-db.sh entrypoint.sh /app/

COPY ./sql /app/sql
RUN chmod +x /app/wait-for-db.sh /app/entrypoint.sh /app/goose
EXPOSE ${EXCHANGER_PORT}
ENTRYPOINT ["/app/entrypoint.sh"]