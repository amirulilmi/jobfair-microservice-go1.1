FROM golang:1.25-alpine AS builder

# Install Git (dan ca-certificates agar HTTPS bisa jalan)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy shared libs first
COPY jobfair-shared-libs/go ./jobfair-shared-libs/go

# Copy ALL service files
COPY jobfair-auth-service/ .

# Update replace directive in go.mod AFTER copying service files
RUN sed -i 's|replace jobfair-shared-libs/go => ../jobfair-shared-libs/go|replace jobfair-shared-libs/go => ./jobfair-shared-libs/go|g' go.mod

# Now download dependencies with corrected go.mod
RUN go mod download

# Build
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
