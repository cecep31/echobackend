FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bin/main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/main .
EXPOSE 8080
CMD ["./main"]