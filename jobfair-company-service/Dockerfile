# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy shared libs
COPY jobfair-shared-libs/go ./jobfair-shared-libs/go/

# Copy go mod files for company service
COPY jobfair-company-service/go.mod jobfair-company-service/go.sum ./jobfair-company-service/

# Download dependencies
WORKDIR /app/jobfair-company-service
RUN go mod download

# Copy source code
COPY jobfair-company-service/ .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/jobfair-company-service/main .

# Expose port
EXPOSE 8080

# Run
CMD ["./main"]
