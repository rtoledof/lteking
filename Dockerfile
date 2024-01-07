FROM golang:1.21.5 AS builder

# first (build) stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /app/auth_service /app/cmd/auth

# final (target) stage

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /app/auth_service ./
CMD [ "./auth_service" ]
