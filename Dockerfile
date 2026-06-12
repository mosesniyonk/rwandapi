# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o rwandapi -ldflags="-s -w" .

# Run stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/rwandapi .

EXPOSE 8080

ENV PORT=8080

CMD ["./rwandapi"]
