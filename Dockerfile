FROM golang:latest
WORKDIR /app/bdiTestService
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN rm Dockerfile
RUN rm README.md
WORKDIR /app/bdiTestService/cmd
RUN go build -o server
RUN rm server.go
EXPOSE 8081
CMD ["./server"]
