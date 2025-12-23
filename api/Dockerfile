FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ffmpeg

WORKDIR /app

# Clone the repository
RUN git clone https://github.com/drnodev/video-stack.git .

# Move into api folder where go.mod exists
WORKDIR /app/api

RUN go mod download

# Build for the CURRENT architecture (important!)
RUN go build -o server

# ---------- Runtime image ----------
FROM alpine:latest

RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

COPY --from=builder /app/api/server .

EXPOSE 8080
CMD ["./server"]
