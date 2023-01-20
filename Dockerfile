FROM golang:latest

WORKDIR /app
COPY ./ /app
RUN go mod download
RUN go install github.com/codegangsta/gin@latest





