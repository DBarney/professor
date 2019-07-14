FROM golang:alpine

RUN apk update && apk add --no-cache git make g++

ADD bin/prof-linux /
WORKDIR /go/.git
ENV GOPATH /go
ENTRYPOINT ["/prof-linux"]
CMD []