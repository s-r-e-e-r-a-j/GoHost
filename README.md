## GoHost
**GoHost** is a powerful HTTP server tool for Debian Linux, written in Go. It allows you to host websites, serve files, and optionally enable downloads for non-web files. With built-in support for **Serveo** and **Cloudflared** tunnels, you can expose your server to the internet and make it accessible from anywhere.

## Features

- Host websites and serve any folder as an HTTP server on Linux.

- Render `.html`, `.htm`, `.css`, and `.js` files as websites.

- When download mode is **off**, files can only be viewed in the browser by clicking on them.

- When download mode is **on**, non-web files can be downloaded by clicking on them.

- Access your server from anywhere using Serveo or Cloudflared tunnels.

- Automatically installs SSH (for Serveo) and Cloudflared if missing.

## Compatibility
- Linux (Debian)

## Installation
1. **Make sure you have Go installed**
   
2. **Clone the GoHost repository:**
   
```bash
git clone https://github.com/s-r-e-e-r-a-j/GoHost.git
```
3. **Navigate to the GoHost folder:**

```bash
cd GoHost
```

4. **Run the installer:**

```bash
sudo bash install.sh
```
 **After installation, you can use `gohost` from anywhere.**

## Usage

```bash
gohost [options]
```

## Options

| Option      | Description                                                   | Default Value     |
|-------------|---------------------------------------------------------------|-------------------|
| `-port`     | Port to run the HTTP server                                   | `8000`            |
| `-path`     | Folder path to serve                                          | Current directory |
| `-tunnel`   | Tunnel method for internet access (`serveo` or `cloudflared`) | None (local only) |
| `-download` | Enable download mode (`true` or `false`)                      | `false`           |

## Examples
1. **Host the current directory on port 8000 (default, view only)**
```bash
gohost
```
3. **Host the current directory on a custom port (view only)**
```bash
gohost -port 8080
```
5. **Host a specific folder on a custom port (view only)**
```bash
gohost -port 8080 -path /home/user/files
```
4. **Host a folder with download mode enabled (non-web files downloadable)**
```bash
gohost -port 8080 -path /home/user/files -download true
```
5. **Expose a folder on the internet**
```bash
# Expose using Serveo 
gohost -port 8080 -path /home/user/files -tunnel serveo

# Expose using Cloudflared 
gohost -port 8080 -path /home/user/files -tunnel cloudflared

# Expose using Cloudflared and enable download mode for non-web files
gohost -port 8080 -path /home/user/files -tunnel cloudflared -download true
```
## Host a Website (HTML/CSS/JS Only)

**Requirement**: Your folder must contain only `.html`, `.css`, and `.js` files
(e.g., `index.html`, `style.css`, `script.js`). No other file types.
```bash
# Default port 8000
gohost -path /home/user/website

# Custom port
gohost -port 8080 -path /home/user/website
```

**Expose the website on the internet:**
```bash
# Serveo
gohost -port 8080 -path /home/user/website -tunnel serveo

# Cloudflared
gohost -port 8080 -path /home/user/website -tunnel cloudflared
```
## Uninstallation 
**Run uninstall.sh script:**
```bash
sudo bash uninstall.sh
```
## License
This project is licensed under the MIT License
