FROM golang:1.21.5 AS builder

# first (build) stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /app/order_service /app/cmd/order

# final (target) stage

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /app/order_service ./
CMD [ "./order_service" ]
