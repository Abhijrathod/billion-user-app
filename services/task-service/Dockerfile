# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY pkg/ ../pkg/
COPY services/task-service/go.mod services/task-service/go.sum* ./
RUN go mod download

# Copy source code
COPY services/task-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 3004

CMD ["./main"]

