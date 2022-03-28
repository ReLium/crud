FROM golang:1.17 AS builder

COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY --from=builder /src/main /app/main
RUN chmod +x /app/main

ENTRYPOINT ["/app/main", "server"]
