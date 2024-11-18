# sudo docker build -t fawazsullialabs/innovelabs-micro-apis:0.0.1 .

# Use the official Golang image to create a build artifact.
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set build environment for cross-platform compatibility
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

COPY .env /root/.env
# Ensure executable permissions
RUN chmod +x main

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]
