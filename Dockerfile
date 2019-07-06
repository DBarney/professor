FROM golang:alpine

RUN apk update && apk add --no-cache git make g++

ADD bin/prof-linux /
WORKDIR /go/.git
ADD www/index.html /go/.git/www
ENV GOPATH /go
ENTRYPOINT ["/prof"]
CMD []