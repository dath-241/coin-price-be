# FROM golang:1.23.2 AS builder
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# RUN go build -o main .
# EXPOSE 8000
# CMD ["./main"]


FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

EXPOSE 8000
CMD ["/app/main"]
