FROM golang:latest
WORKDIR /app/bdiTestService
COPY go.mod ./
RUN go mod download
COPY . .
WORKDIR /app/bdiTestService/cmd
RUN go build -o server
RUN rm server.go
WORKDIR /app/bdiTestService
RUN rm -rf app && rm Dockerfile && rm go.mod && rm go.sum && rm README.md
WORKDIR /app/bdiTestService/cmd
EXPOSE 8081
CMD ["./server"]
