# Stage 1: Dependencies
FROM golang:1.22 AS deps

WORKDIR /hello-api
ADD *.mod *.sum ./
RUN go mod download

# Stage 2: Development and Build stage
FROM deps as dev

# Add source code
ADD . .

# Expose the port
EXPOSE 8080

# Build arguments for LDFLAGS
ARG GO_VERSION=1.22
ARG TAG
ARG HASH
ARG DATE

# Set up LDFLAGS with provided build arguments
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags \
    "-w -X main.docker=true \
    -X github.com/kazakh-in-nz/hello-api/handlers.hash=${HASH} \
    -X github.com/kazakh-in-nz/hello-api/handlers.tag=${TAG} \
    -X github.com/kazakh-in-nz/hello-api/handlers.date=${DATE}" \
    -o api cmd/main.go

# Command to run the application in development
CMD ["/hello-api/api"]

# Stage 3: Production
FROM scratch as prod

WORKDIR /
EXPOSE 8080

# Copy the built binary from the dev stage
COPY --from=dev /hello-api/api /

# Run the binary
CMD ["/api"]
