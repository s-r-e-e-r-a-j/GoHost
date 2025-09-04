// Developer: Sreeraj 
// GitHub: https://github.com/s-r-e-e-r-a-j
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"path/filepath"
	"os/signal"
    "syscall"
    "time"
)

// checkInstall installs a command if not found
func checkInstall(command string, installCmd []string) {
	_, err := exec.LookPath(command)
	if err != nil {
		fmt.Printf("%s not found. Installing...\n", command)
		cmd := exec.Command(installCmd[0], installCmd[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to install %s: %v", command, err)
		}
	}
}

// startServeo starts a Serveo tunnel
func startServeo(port int) *exec.Cmd {
	fmt.Println("[*] Starting Serveo tunnel...")
	cmd := exec.Command("ssh", "-R", fmt.Sprintf("80:localhost:%d", port), "serveo.net")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Serveo: %v", err)
	}
	return cmd
}

// startCloudflared starts Cloudflared tunnel
func startCloudflared(port int) *exec.Cmd {
	fmt.Println("[*] Starting Cloudflared tunnel...")
	cmd := exec.Command("cloudflared", "tunnel", "--url", fmt.Sprintf("http://localhost:%d", port))

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start cloudflared: %v", err)
	}

	// Regex to match trycloudflare.com URL
	re := regexp.MustCompile(`https://[^\s]+\.trycloudflare\.com`)
	var lastURL string


	scanPipe := func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			if match := re.FindString(line); match != "" {
				lastURL = match
			}
		}
	}

	go scanPipe(bufio.NewScanner(stdout))
	go scanPipe(bufio.NewScanner(stderr))

	// Wait until a URL is found
	for lastURL == "" {
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("[Cloudflared URL] %s\n", lastURL)
	return cmd
}

// fileHandler
func fileHandler(path string, forceDownload bool) http.Handler {
	webExts := map[string]bool{
		".html": true,
		".htm":  true,
		".css":  true,
		".js":   true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := path + r.URL.Path
		stat, err := os.Stat(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if stat.IsDir() {
			indexPath := filePath + "/index.html"
			if _, err := os.Stat(indexPath); err == nil && !forceDownload {
				http.ServeFile(w, r, indexPath)
				return
			}
			http.FileServer(http.Dir(path)).ServeHTTP(w, r)
			return
		}

		if forceDownload {
		    ext := filepath.Ext(filePath)
            if !webExts[ext] {
                 w.Header().Set("Content-Disposition", "attachment; filename="+stat.Name())
                 w.Header().Set("Content-Type", "application/octet-stream")
            }			
		}

		http.ServeFile(w, r, filePath)
	})
}

func main() {
	if runtime.GOOS != "linux" {
		log.Fatal("This server is only supported on Linux.")
	}

	port := flag.Int("port", 8000, "Port to run HTTP server")
	tunnel := flag.String("tunnel", "", "Port forwarding option: serveo or cloudflared")
	path := flag.String("path", "", "Folder path to serve (optional, default current directory)")
	download := flag.Bool("download", false, "Enable download mode (true/false)")
	flag.Parse()

	// Determine folder to serve
	servePath := *path
	if servePath == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		servePath = dir
	} else if _, err := os.Stat(servePath); os.IsNotExist(err) {
		log.Fatalf("Folder %s does not exist", servePath)
	}

	// Install SSH if using Serveo
	if *tunnel == "serveo" {
		checkInstall("ssh", []string{"sudo", "apt", "install", "-y", "openssh-client"})
	}

	// Install Cloudflared if using Cloudflared
	if *tunnel == "cloudflared" {
		checkInstall("cloudflared", []string{"sudo", "wget", "-qO", "/usr/local/bin/cloudflared", "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64"})
		
		cloudflaredPath := "/usr/local/bin/cloudflared"
        if info, err := os.Stat(cloudflaredPath); err == nil {
            mode := info.Mode()
            if mode&0111 == 0 { // if not executable
                fmt.Println("[*] Setting execute permission for Cloudflared...")
                cmd := exec.Command("sudo", "chmod", "+x", cloudflaredPath)
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr
                if err := cmd.Run(); err != nil {
                    log.Fatalf("Failed to set execute permission for Cloudflared: %v", err)
                 }
             }
        }
	}

	// Serve files
	fmt.Printf("[*] Serving %s on http://localhost:%d\n", servePath, *port)
	http.Handle("/", fileHandler(servePath, *download))

    var serveoCmd *exec.Cmd
    var cloudflaredCmd *exec.Cmd

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigChan
        fmt.Println("\n[*] Caught interrupt signal, cleaning up...")

        if serveoCmd != nil && serveoCmd.Process != nil {
            fmt.Println("[*] Killing Serveo tunnel...")
            serveoCmd.Process.Kill()
        }

        if cloudflaredCmd != nil && cloudflaredCmd.Process != nil {
            fmt.Println("[*] Killing Cloudflared tunnel...")
            cloudflaredCmd.Process.Kill()
        }

        os.Exit(0)
     }()
		
    // Start tunnel if specified
    if *tunnel == "serveo" {
	     serveoCmd = startServeo(*port)
    } else if *tunnel == "cloudflared" {
	           cloudflaredCmd = startCloudflared(*port)
    }
	
	// Start HTTP server
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
