FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
RUN apk --no-cache add ca-certificates
EXPOSE 8080
CMD ["./server"]