# Use an official Golang base image 
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container for the application build
WORKDIR /app

# Copy go modules files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o main ./main.go

# Use a lighweight image for running the app
FROM alpine:latest

# Set the working directory in the lightweight image
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file
COPY .env .
COPY conf/. ./conf/

# Expose the port on which the app will listen
EXPOSE 8080

# Run the Go app
CMD ["./main"]