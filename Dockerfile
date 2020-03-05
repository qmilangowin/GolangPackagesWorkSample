FROM golang:alpine AS builder
# RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"
WORKDIR /app/bdiTestService
COPY . .
RUN go mod download
WORKDIR /app/bdiTestService/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bdi_test_service
RUN rm server.go
WORKDIR /app/bdiTestService/
RUN rm -rf app && rm -rf tests && rm Dockerfile && rm go.mod && rm go.sum && rm README.md
#FROM golang:alpine
#RUN rm -rf app && rm Dockerfile && rm go.mod && rm go.sum && rm README.md
#COPY --from=builder app/bdiTestService/cmd /usr/local/bin
WORKDIR /app/bdiTestService/cmd
EXPOSE 8081
#ENTRYPOINT ["server"]
CMD ["./bdi_test_service"]
