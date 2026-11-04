
# Stage 1: Build the Go application
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Copy dependency files first for better Docker layer caching
COPY go.mod go.sum ./

RUN go mod download

# Copy application source code
COPY . .

# Run automated tests during image build
RUN go test ./...

# Compile the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o employee-service .

# Stage 2: Create the runtime image
FROM alpine:3.21

WORKDIR /app

# Install CA certificates for future HTTPS communication
RUN apk --no-cache add ca-certificates

# Copy only the compiled binary
COPY --from=builder /app/employee-service .

# Run as a non-root user
USER 10001:10001

EXPOSE 8080

CMD ["./employee-service"]
