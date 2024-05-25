# Use the official Golang image as the build stage
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and go sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app and place the output binary in the /app directory
RUN CGO_ENABLED=0 GOOS=linux  go build -o /app/cmd/api/main ./cmd/api

# Start a new stage from scratch
FROM alpine:latest  

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/cmd/api/main .

RUN chmod +x main

# Command to run the executable
ENTRYPOINT ["./main"]

