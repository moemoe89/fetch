# Stage 1: Build the Go binary
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o fetch ./cmd/main.go

# Stage 2: Copy the binary to a minimal Alpine image
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/fetch /app/
CMD ["./fetch", "$@"]
