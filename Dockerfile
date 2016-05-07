FROM ubuntu:14.04
MAINTAINER Andrey Gorelov <ekonomizer@gmail.com>

RUN apt-get update && apt-get install -y golang