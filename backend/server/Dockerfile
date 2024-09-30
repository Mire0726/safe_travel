# Start from the official Go image
FROM golang:1.21

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install sqlboiler and its MySQL driver
RUN go install github.com/volatiletech/sqlboiler/v4@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest

# Build the application
# Assuming main.go is in the server directory
RUN go build -o main ./server

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]