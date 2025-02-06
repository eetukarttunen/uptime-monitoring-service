# Golang image
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod tidy
COPY . .
RUN go build -o uptime-monitoring-service .

# Another small image to run the application
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/uptime-monitoring-service .
CMD ["./uptime-monitoring-service"]
