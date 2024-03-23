# =============================================================================
# Build phase
# =============================================================================

# Define our base image from the official Golang image
FROM golang:latest AS builder

# Create and/or change our directory to /build
WORKDIR /build

# Copy go mod & go.sum files so we can verify they haven't been tampered with
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy everything from our root into /app
COPY . .

# Create and/or change our directory to /cmd
WORKDIR /build/cmd

# Create the binary for the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-binary

# =============================================================================
# Production phase
# =============================================================================

# Start a new stage from scratch so that our final image is way smaller
# We use alpine instead of stratch since we need OS packages
FROM alpine:latest

# Create and/or change our directory to /app
WORKDIR /app

# Copy the binary from the build stage to the final stage
COPY --from=builder /build/cmd/app-binary .

# Expose Port
EXPOSE 8080

# Run the app
CMD [ "/app/app-binary"]