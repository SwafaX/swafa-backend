FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Use a minimal Alpine image for the final stage
FROM alpine:latest

COPY --from=builder /app/main .

COPY .env .env

EXPOSE 8000

CMD ["./main"]
