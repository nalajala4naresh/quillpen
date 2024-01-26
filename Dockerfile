# Stage 1: Build the Go application
FROM golang:1.21-alpine as builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application with CGO disabled and as a statically-linked executable
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o quillpen

# Stage 2: Create the minimal runtime image
FROM scratch

WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /app/quillpen /app/quillpen
COPY ./sf-class2-root.crt  private_key.pem  certificate.pem /app/


# Set the command to run the application
CMD ["/app/quillpen"]