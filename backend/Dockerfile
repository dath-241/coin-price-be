FROM golang:1.23.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go build -o main .
EXPOSE 8000
CMD ["./main"]