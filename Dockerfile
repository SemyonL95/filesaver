FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /app
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
ADD . .

# Build the application
RUN go build -o main ./app/.

# Move to /dist directory as the place for resulting binary folder
WORKDIR /bin

# Copy binary from build to main folder
RUN cp /app/main .

# Build a small image
FROM scratch

COPY --from=builder /app/main /

# Command to run
ENTRYPOINT ["/main"]

EXPOSE 8080

VOLUME [ "/storage", "/tmp" ]