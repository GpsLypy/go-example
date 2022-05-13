# This file should be run in GOPATH
#FROM docker.17usoft.com/library/golang:1.7
FROM docker.17usoft.com/library/centos:7.2

ENV TZ="Asia/Shanghai"

RUN mkdir /usr/local/mover
ADD bin/mover.tar.gz /usr/local/mover

WORKDIR /usr/local/mover/bin
USER root

CMD ["./mover"]
