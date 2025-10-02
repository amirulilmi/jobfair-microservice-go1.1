FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY jobfair-api-gateway/go.mod jobfair-api-gateway/go.sum ./
RUN go mod download

# Copy source code
COPY jobfair-api-gateway/ .

# Build
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
