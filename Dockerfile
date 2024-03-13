FROM golang:1.22.1-alpine3.19 AS builder

COPY *.go go.mod go.sum /go/src/

WORKDIR /go/src

RUN go mod tidy && \
  GOOS=linux GOARCH=arm go build -o myip

FROM alpine:3.19.1

WORKDIR /

COPY --from=builder /go/src/myip /

ENTRYPOINT [ "/myip" ]
