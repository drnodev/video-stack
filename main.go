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
	start := time.Now()
	logPrefix := fmt.Sprintf("[%s]", start.Format(time.RFC3339))

	fmt.Println(logPrefix, "Incoming request")
	fmt.Println(logPrefix, "Method:", r.Method)
	fmt.Println(logPrefix, "RemoteAddr:", r.RemoteAddr)
	fmt.Println(logPrefix, "Content-Type:", r.Header.Get("Content-Type"))
	fmt.Println(logPrefix, "Content-Length:", r.ContentLength)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Println(logPrefix, "ParseMultipartForm error:", err)
		http.Error(w, "File too large or invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		fmt.Println(logPrefix, "FormFile error:", err)
		http.Error(w, "Missing video field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println(logPrefix, "Received file:")
	fmt.Println(logPrefix, "Filename:", header.Filename)
	fmt.Println(logPrefix, "Size:", header.Size)
	fmt.Println(logPrefix, "MIME:", header.Header.Get("Content-Type"))

	tmpDir := os.TempDir()
	inputPath := filepath.Join(tmpDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename))
	outputPath := filepath.Join(tmpDir, "noaudio_"+filepath.Base(inputPath))

	fmt.Println(logPrefix, "Input path:", inputPath)
	fmt.Println(logPrefix, "Output path:", outputPath)

	out, err := os.Create(inputPath)
	if err != nil {
		fmt.Println(logPrefix, "Create input file error:", err)
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(out, file); err != nil {
		fmt.Println(logPrefix, "File copy error:", err)
		out.Close()
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	out.Close()

	// 👇 IMPORTANTE: capturar stdout + stderr de ffmpeg
	cmd := exec.Command(
		"ffmpeg",
		"-hide_banner",
		"-loglevel", "error", // cambia a "info" o "verbose" si quieres MÁS
		"-y",
		"-i", inputPath,
		"-c:v", "copy",
		"-an",
		outputPath,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(logPrefix, "FFmpeg command failed")
		fmt.Println(logPrefix, "FFmpeg error:", err)
		fmt.Println(logPrefix, "FFmpeg output:\n", string(output))

		http.Error(w, "FFmpeg failed", http.StatusInternalServerError)
		return
	}

	fmt.Println(logPrefix, "FFmpeg completed successfully")
	fmt.Println(logPrefix, "FFmpeg output:\n", string(output))

	defer os.Remove(inputPath)
	defer os.Remove(outputPath)

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=no_audio_"+header.Filename,
	)

	http.ServeFile(w, r, outputPath)

	fmt.Println(logPrefix, "Request completed in", time.Since(start))
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodGet {
		fmt.Fprint(w, "😊")
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/remove-audio", removeAudioHandler)
	fmt.Println("🚀 Go Video API running on :9090")
	http.ListenAndServe(":9090", nil)
}
