FROM golang:1.8
MAINTAINER g-hyoga <hyoga0216@gmail.com>

# RUN apt-get -y install mysql-server

RUN mkdir /go/src/github.com \ 
	&& mkdir /go/src/github.com/g-hyoga
RUN cd /go/src/github.com/g-hyoga && git clone https://github.com/g-hyoga/kyuko.git

WORKDIR /go/src/github.com/g-hyoga/kyuko/go
RUN make setup

