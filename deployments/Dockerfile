FROM golang:1.22.3-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Build app
COPY . .
RUN go build -o server ./cmd/blog_engine/main.go

# Build final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server ./

CMD ["/app/server"]
