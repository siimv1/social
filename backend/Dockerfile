# Use the official Golang image as a base
FROM golang:1.22.0 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code, including the backend folder
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM golang:1.22.0

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Ensure the migrations are correctly copied to match the relative path used in main.go
COPY --from=builder /app/backend/pkg/db/migrations /app/backend/pkg/db/migrations

# Set the working directory
WORKDIR /app

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]
