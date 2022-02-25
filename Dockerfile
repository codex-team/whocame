FROM golang:latest
ARG BUILD_DIRECTORY=/build

RUN apt update
RUN apt install -y libpcap-dev

WORKDIR /app
RUN mkdir $BUILD_DIRECTORY
COPY go.mod .
COPY go.sum .


RUN go mod download

ADD . $BUILD_DIRECTORY