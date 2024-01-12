FROM golang:1.21.5 AS builder

# first (build) stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /app/rider_service /app/cmd/rider

# final (target) stage

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /app/rider_service ./
CMD [ "./rider_service" ]
