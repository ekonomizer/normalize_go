FROM ubuntu:16.04
MAINTAINER Andrey Gorelov <ekonomizer@gmail.com>


RUN mkdir project
ENV GOPATH /project
RUN apt-get update && apt-get install -y git && apt-get install -y golang
RUN apt-get install unzip