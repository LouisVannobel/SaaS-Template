# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files from backend directory
COPY ./backend/go.mod ./backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY ./backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest AS runner

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy any necessary config files
COPY --from=builder /app/db/migrations ./migrations

# Create a non-root user and switch to it
RUN adduser -D appuser
RUN chown -R appuser:appuser /app
USER appuser

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]