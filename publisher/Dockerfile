FROM golang:1.17 AS builder

WORKDIR /build
COPY . .
RUN go get -d -v ./...

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o publisher .

FROM alpine:3.13 AS certificates
RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /app
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/publisher .

EXPOSE 9001

CMD ["./publisher"]