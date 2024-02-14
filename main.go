package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing URL parameter", http.StatusBadRequest)
			return
		}

		cmd := exec.Command("ffmpeg", "-i", url, "-c", "copy", "-f", "mp4", "-y", "pipe:1")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.StdoutPipe()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create stdout pipe: %v", err), http.StatusInternalServerError)
			return
		}

		if err := cmd.Start(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to start command: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "video/mp4")
		if _, err := io.Copy(w, out); err != nil {
			http.Error(w, fmt.Sprintf("Failed to copy output: %v", err), http.StatusInternalServerError)
			return
		}

		if err := cmd.Wait(); err != nil {
			http.Error(w, fmt.Sprintf("Command execution failed: %v: %s", err, stderr.String()), http.StatusInternalServerError)
			return
		}
	})
	
	http.ListenAndServe(":8080", nil)
}
