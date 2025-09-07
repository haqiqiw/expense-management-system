# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git build-base librdkafka-dev pkgconf

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o expense-approved-consumer ./cmd/expense-approved-consumer

# Runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/expense-approved-consumer .

CMD ["./expense-approved-consumer"]