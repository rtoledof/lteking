FROM golang:1.21.5 AS builder

# first (build) stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /app/wallet_service /app/main

# final (target) stage

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /app/wallet_service ./
CMD [ "./wallet_service" ]
