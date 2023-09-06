FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src

COPY go.mod .

RUN go mod download

COPY . .

RUN go build src/main.go

RUN apk update && apk upgrade && apk add --no-cache ca-certificates

RUN mkdir -p /etc/ssl/certs/ && update-ca-certificates

FROM scratch

COPY --from=builder /go/src .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT  ["./main"]