# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY jobfair-user-profile-service/go.mod jobfair-user-profile-service/go.sum ./
RUN go mod download

# Copy source code
COPY jobfair-user-profile-service/ .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Create upload directory
RUN mkdir -p ./uploads/cv

# Expose port
EXPOSE 8083

# Run
CMD ["./main"]
