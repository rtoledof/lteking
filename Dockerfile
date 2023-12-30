FROM golang:1.21.4 AS builder

# first (build) stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /app/cubawheeler /app/cmd/cubawheeler

# final (target) stage

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /app/cubawheeler ./
CMD [ "./cubawheeler" ]
