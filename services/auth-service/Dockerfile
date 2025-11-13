# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy workspace file
COPY go.work ./

# Copy all pkg modules
COPY pkg/ ./pkg/

# Copy service go.mod
COPY services/auth-service/go.mod ./services/auth-service/
WORKDIR /app/services/auth-service
RUN go mod download

# Copy service source
COPY services/auth-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/auth-service/main .

EXPOSE 3001

CMD ["./main"]

