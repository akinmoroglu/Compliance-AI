FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/

FROM alpine:3.19

RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/prompts ./prompts

EXPOSE 8080

CMD ["./server"]
