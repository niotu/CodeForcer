# Stage 1: Setup Node.js environment
FROM node:20-alpine as frontend

# Set the working directory for the Node.js application
WORKDIR /app

# Copy the Node.js application code
COPY web/front/ ./web/front/
RUN npm install --prefix ./web/front

# Stage 2: Build the Go application
FROM golang:1.21

# Set the working directory for the Go application
WORKDIR /app

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go application code
COPY ./ ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./backend

COPY --from=frontend /app/ ./

# Ensure the server binary is executable
RUN chmod +x ./server

# Make the script executable
RUN chmod +x ./start.sh

# Expose port 8080
EXPOSE 8080
EXPOSE 80

# Set the CMD to run the script
CMD ["./start.sh"]
