# Use the official Golang image as the build stage
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp ./example
RUN ls -la /app  # Check if the binary is created

# Use an Ubuntu base image for the final stage
FROM ubuntu:22.04

# Set the working directory inside the container
WORKDIR /app

# Install necessary libraries
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the built Go application from the builder stage
COPY --from=builder /app/myapp .
RUN ls -la /app  # Check if the binary is copied
RUN chmod +x /app/myapp  # Ensure the binary is executable

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]
