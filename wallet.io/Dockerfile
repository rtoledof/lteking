FROM golang:latest AS builder

# first (build) stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /app/wallet_service /app

# final (target) stage

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /app/wallet_service ./
CMD [ "./wallet_service" ]
