FROM ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.21.3.linux-amd64.tar.gz

WORKDIR /go-http3-proxy
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o server ./cmd/server

CMD ./server
