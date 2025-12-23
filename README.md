# Video Audio Remover API (Go + FFmpeg)

This project provides a lightweight **Go-based REST API** that removes the audio track from a video using **FFmpeg**.  
It is designed to run as a **microservice**, ideal for automation tools like **n8n** or any system that needs video processing via HTTP.

The entire stack is containerized using **Docker Compose**.

---

## 📁 Project Structure

```video-stack/
├── api/
│ ├── main.go # Go HTTP server
│ ├── go.mod # Go module definition
│ └── Dockerfile # API container with FFmpeg
└── docker-compose.yml # Orchestrates API
```

---

## 🚀 Prerequisites

Make sure you have the following installed:

- Docker
- Docker Compose (v2+ recommended)

Verify installation:

```bash
docker --version
docker compose version
```

---

## 🏗️ Setup

Build and run:
```bash
docker compose up --build
```

The API will be available at `http://localhost:108080`.

---

## 📝 API Endpoints

### Remove Audio

- **Endpoint:** `POST /remove-audio`
- **Request:**
  ```http
  POST /remove-audio HTTP/1.1
  Content-Type: multipart/form-data

  --boundary
  Content-Disposition: form-data; name="video"; filename="your-video.mp4"
  Content-Type: video/mp4

  <video-file>
  --boundary--
  ```

- **Response:**
  ```http
  HTTP/1.1 200 OK
  Content-Type: video/mp4
  Content-Disposition: attachment; filename="no_audio_your-video.mp4"

  <video-file-without-audio>
  ```

---

## 📊 Performance

- The API is optimized for performance with:
  - Efficient file handling
  - Minimal memory usage
  - Fast processing
