# Use an official Go runtime as a parent image
FROM golang:1.22

# Set the working directory to /app
WORKDIR /app

# Copy your Go code into the container at /app
COPY . /app

# Build your Go binary
RUN go build -o cmd/main cmd/main.go

# Make port 8080 available to the world outside this container
EXPOSE 8080

# Run your Go binary when the container launches
CMD ["/app/cmd/main"]