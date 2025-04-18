# Build stage
FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sustainability-tracker-activity-service ./cmd

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sustainability-tracker-activity-service .
EXPOSE 8081
CMD ["./sustainability-tracker-activity-service"]