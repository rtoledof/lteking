FROM golang:latest as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o order cmd/order/main.go
EXPOSE  3000
CMD ["./order"]