# Start from the base Go image
FROM golang:1.17 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.* ./

# Download all dependencies
RUN go mod download

# Copy the entire source code from the current directory to the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

####### Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside
# This assumes your Go app runs on port 8080, change if necessary
EXPOSE 8080

# Run the binary
CMD ["./main"]
