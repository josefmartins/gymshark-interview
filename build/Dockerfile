FROM golang:1.24.2 AS builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./product-service ./cmd/main.go

RUN chmod +x ./product-service

FROM debian:12.10-slim

COPY --from=builder /app/product-service /product-service
COPY --from=builder /app/database/migrations/*.sql /database/migrations/

EXPOSE 8080

CMD ["/product-service"]