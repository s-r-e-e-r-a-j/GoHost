#!/bin/bash

if [[ "$EUID" -ne 0 ]]; then
    echo "Please run as root or with sudo."
    exit 1
fi

# Check if gohost.go exists
if [[ ! -f "gohost.go" ]]; then
    echo "gohost.go not found in this directory. Please place it here."
    exit 1
fi

# Build the Go program
echo "[*] Building gohost from ghost.go..."
go build -o gohost gohost.go

# Move executable to /usr/local/bin
echo "[*] Installing gohost as a system-wide command..."
sudo mv gohost /usr/local/bin/
sudo chmod +x /usr/local/bin/gohost

echo "[*] Installation complete!"
echo "You can now run 'gohost' from anywhere:"
echo "Example: gohost -port 8080 -path /home/user/website -tunnel serveo"
