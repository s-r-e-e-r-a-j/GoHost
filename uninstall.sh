#!/bin/bash

# Check if running as root
if [[ "$EUID" -ne 0 ]]; then
    echo "Please run as root or with sudo."
    exit 1
fi

# Check if gohost exists in /usr/local/bin
if [[ -f "/usr/local/bin/gohost" ]]; then
    echo "[*] Removing gohost from /usr/local/bin..."
    sudo rm -f /usr/local/bin/gohost
    echo "[*] gohost has been uninstalled."
else
    echo "[!] gohost is not installed in /usr/local/bin."
fi
