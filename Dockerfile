# Stage 1: Build the Go application
FROM golang:1.21 AS builder

# Set the working directory for the Go application
WORKDIR /app

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go application code except web/ directory
COPY backend ./backend

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./backend

# Stage 2: Setup Node.js environment and combine both applications
FROM node:20-alpine

# Set the working directory for the combined application
WORKDIR /app

# Copy the Node.js application code
COPY web/ ./web
COPY db/ ./db
RUN npm install --prefix ./web/front

# Copy the built Go server binary from the first stage
COPY --from=builder /app/server ./server

# Copy the start script and ensure it is executable
COPY start.sh ./start.sh
RUN chmod +x ./server ./start.sh

# Copy the rest of the frontend code
COPY credentials.json ./

# Expose ports
EXPOSE 8080
EXPOSE 80

# Set the CMD to run the script
CMD ["./start.sh"]
