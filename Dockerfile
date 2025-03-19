FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go install github.com/google/wire/cmd/wire@latest
RUN wire internal/di/wire.go
RUN go build -o bin/main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/main .
EXPOSE 8080
CMD ["./main"]