# Use the Go 1.21 base image
FROM golang:1.21 as base

# Install Node.js
RUN apt-get update && apt-get install -y curl \
  && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
  && apt-get install -y nodejs \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Set the working directory for the application
WORKDIR /app

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go application code except web/ directory
COPY ./ ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./backend

# Setup Node.js environment
WORKDIR /app/web/front
RUN npm install

# Move back to the app directory
WORKDIR /app

# Ensure the server binary is executable
RUN chmod +x ./server

# Copy the start script and ensure it is executable
RUN chmod +x ./start.sh

# Expose ports
EXPOSE 8080
EXPOSE 80

# Set the CMD to run the script
CMD ["./start.sh"]
