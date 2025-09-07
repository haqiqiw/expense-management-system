# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags musl -o mock-payment-api ./dev/payment-api

# Runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/mock-payment-api .

EXPOSE 9500

CMD ["./mock-payment-api"]