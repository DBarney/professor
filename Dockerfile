FROM golang:1.12-stretch

RUN apt-get install git make g++

ADD bin/prof-linux /
WORKDIR /go/.git
ENV GOPATH /go
ENTRYPOINT ["/prof-linux"]
CMD []