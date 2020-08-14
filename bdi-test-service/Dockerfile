FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#set directory to map datasets
RUN mkdir /home/data
# create a build directory
WORKDIR /build

# Copy and download dependency using go mod
COPY . .
RUN go mod download

# Build the service
RUN go build -o bdi-test-service .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/bdi-test-service .

# Build a small image
FROM scratch

COPY --from=builder /dist/bdi-test-service /
#EXPOSE 8081
# Command to run
CMD ["/bdi-test-service"]
