# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy shared libs
COPY jobfair-shared-libs/go ./jobfair-shared-libs/go/

# Copy go mod files for job service
COPY jobfair-job-service/go.mod jobfair-job-service/go.sum ./jobfair-job-service/

# Download dependencies
WORKDIR /app/jobfair-job-service
RUN go mod download

# Copy source code
COPY jobfair-job-service/ .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/jobfair-job-service/main .

# Expose port
EXPOSE 8082

# Run
CMD ["./main"]
