FROM golang:1.8-alpine

RUN apk update && apk add bash git

# Add the files & install
COPY . /go/src/github.com/scristofari/image-server/part1
WORKDIR /go/src/github.com/scristofari/image-server/part1
RUN go install

CMD image-server