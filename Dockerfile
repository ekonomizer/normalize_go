FROM ubuntu:14.04
MAINTAINER Andrey Gorelov <ekonomizer@gmail.com>

RUN apt-get update && apt-get install -y golang && apt-get install -y git
RUN mkdir projects && cd projects && mkdir go && cd go && git init && git remote add origin https://github.com/ekonomizer/normalize_go.git && git fetch && git checkout master
