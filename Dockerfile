# Stage 1: Build the Go application
FROM golang:1.21 AS builder

# Set the working directory
WORKDIR /backend

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /backend/backend

# Copy the rest of the application code
COPY backend credentials.json  ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./server.go

# Stage 2: Setup Node.js environment
FROM node:20-alpine

WORKDIR /frontend

COPY web/front/ ./

RUN npm install

# Copy the built Go application from the builder stage
COPY --from=builder /backend/backend/server ./server

# Ensure the server binary is executable
RUN chmod +x ./server

# Copy the script into the Docker image
COPY start.sh ./start.sh

# Make the script executable
RUN chmod +x ./start.sh

# Expose port 8080
EXPOSE 8080
EXPOSE 80

# Set the CMD to run the script
CMD ["./start.sh"]
