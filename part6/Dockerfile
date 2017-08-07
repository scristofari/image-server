FROM golang:1.8-alpine as builder
COPY . /go/src/github.com/scristofari/image-server/part6
WORKDIR /go/src/github.com/scristofari/image-server/part6/cmd/server
RUN go build -o image-server

FROM alpine:latest
WORKDIR /go/app
COPY --from=builder /go/src/github.com/scristofari/image-server/part6 .
WORKDIR /go/app/cmd/server
RUN adduser -D imageuser
USER imageuser
CMD ./image-server
