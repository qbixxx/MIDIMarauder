#!/bin/bash

# Ensure Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go before proceeding."
    exit 1
fi

# Check if pkg-config is installed
if ! command -v pkg-config &> /dev/null; then
    echo "pkg-config is not installed. Installing..."
    
    # Install pkg-config
    # For Debian/Ubuntu, adjust for your distribution
    sudo apt update
    sudo apt install pkg-config
    
    # Check if installation was successful
    if ! command -v pkg-config &> /dev/null; then
        echo "Failed to install pkg-config. Please install it manually."
        exit 1
    fi
fi

# Initialize a Go module
go mod init midimarauder

# Run 'go mod tidy' to clean up dependencies
go mod tidy

# Install libusb (assuming it's available as a package)
# For Debian/Ubuntu, adjust for your distribution
sudo apt install libusb-1.0-0-dev

# Build the Go code
go build midimarauder.go

# Make the executable accessible system-wide
sudo cp "./midimarauder" "/usr/local/bin/"
sudo chmod u+s "/usr/local/bin/midimarauder"

# Check if the program is now available
if command -v midimarauder &> /dev/null; then
    echo "Installation successful. You can now run 'midimarauder' from the terminal."
else
    echo "Installation failed. Check for errors above."
fi
