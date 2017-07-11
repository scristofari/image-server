FROM golang:1.8-alpine

RUN apk update && apk add bash git

# never root
RUN adduser -D imageuser
USER imageuser

# add dependency management tool
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/scristofari/image-server
WORKDIR /go/src/github.com/scristofari/image-server
RUN go install

CMD image-server