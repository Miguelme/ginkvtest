# Stage 1: Build the application
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Set environment variables for static build
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . .

# Build the application binary
RUN go build -o main cmd/main.go

# Stage 2: Run the application using distroless
FROM gcr.io/distroless/static:nonroot

# Set the working directory
WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the executable
CMD ["/main"]
