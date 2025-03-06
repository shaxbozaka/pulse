# Multi-stage build for smaller final image
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

# Build the application
RUN go build -o pulse-server server/server.go

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/pulse-server .

EXPOSE 50051
ENTRYPOINT ["./pulse-server"]
