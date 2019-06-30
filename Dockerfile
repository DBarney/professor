FROM golang:alpine

RUN apk update && apk add --no-cache git make g++

ADD bin/prof /
WORKDIR /go/.git
ENV GOPATH /go
ENTRYPOINT ["/prof"]
CMD []