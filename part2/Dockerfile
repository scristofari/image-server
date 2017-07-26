FROM golang:1.8-alpine

RUN apk update && apk add bash git

# add files
COPY . /go/src/github.com/scristofari/image-server/part2
WORKDIR /go/src/github.com/scristofari/image-server/part2

# add dependency management tool
RUN go get -u github.com/golang/dep/cmd/dep

# install and run
RUN go install
CMD image-server