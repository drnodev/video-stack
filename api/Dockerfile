FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ffmpeg

WORKDIR /app

RUN git clone https://github.com/drnodev/video-stack.git .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

FROM alpine:latest
RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]