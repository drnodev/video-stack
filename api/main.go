package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const maxUploadSize = 500 << 20

func removeAudioHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Missing video field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpDir := os.TempDir()
	inputPath := filepath.Join(tmpDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename))
	outputPath := filepath.Join(tmpDir, "noaudio_"+filepath.Base(inputPath))

	out, _ := os.Create(inputPath)
	io.Copy(out, file)
	out.Close()

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", inputPath,
		"-c:v", "copy",
		"-an",
		outputPath,
	)

	if err := cmd.Run(); err != nil {
		http.Error(w, "FFmpeg failed", http.StatusInternalServerError)
		fmt.Println("FFmpeg failed", err)

		return
	}

	defer os.Remove(inputPath)
	defer os.Remove(outputPath)

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", "attachment; filename=no_audio_"+header.Filename)
	http.ServeFile(w, r, outputPath)
}

func main() {
	http.HandleFunc("/remove-audio", removeAudioHandler)
	fmt.Println("🚀 Go Video API running on :8080")
	http.ListenAndServe(":8080", nil)
}
