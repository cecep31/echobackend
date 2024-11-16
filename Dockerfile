from golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bin/main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/main .
EXPOSE 1323
CMD ["./main"]