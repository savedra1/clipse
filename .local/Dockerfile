# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Install dependencies for GLFW and GUI applications
RUN apt-get update && apt-get install -y \
    libgl1-mesa-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxxf86vm-dev \
    libasound2-dev \
    libudev-dev \
    libdbus-1-dev \
    libx11-dev

# Download and install Fyne dependencies
RUN go get fyne.io/fyne/v2

# Expose the port the application will run on
EXPOSE 8080

# Command to run the Go application
CMD ["go", "run", "main.go"]
