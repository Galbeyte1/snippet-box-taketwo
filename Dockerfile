# Set the base image, alias as the builder stage
FROM golang:1.23 AS builder

# Create an app directory
RUN mkdir /app

# Add all the contents of the current directory to this app directory
ADD . /app

# Specify the working directory as this app directory
WORKDIR /app

# Build the compiled binary of the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/web

# Define the second stage of our multi-stage dockerfile
FROM alpine:latest AS production

# Install certificates (important for HTTPS inside container)
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/ui ./ui

# Specify the command to execute this app compiled binary
CMD ["./app"]

# Benefits of multi stage
# reduced resource usage from the first stage not necessary in second stage

# docker run snippetbox-server  
# or 
# docker run -it -p 8080:8080 snippetbox-server