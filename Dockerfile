FROM golang:1.5.2
MAINTAINER LeReunionais

RUN apt-get update
RUN apt-get install libzmq3 -y
RUN apt-get install libzmq3-dev -y
RUN ldconfig
RUN apt-get install pkg-config -y

COPY . /go/src/github.com/LeReunionais/looper
WORKDIR /go/src/github.com/LeReunionais/looper

RUN go get -d -v
RUN go install -v

EXPOSE 6000

CMD looper
