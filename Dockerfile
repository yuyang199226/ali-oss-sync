FROM ubuntu:16.04

RUN  apt-get update && apt-get -y upgrade
RUN apt-get install libkrb5-dev krb5-user
RUN mkdir -p /go/bin && \
    wget -nv https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz && \
    tar -xvf go1.13.8.linux-amd64.tar.gz -C /usr/local/ && \
    rm go1.13.8.linux-amd64.tar.gz
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

RUN wget http://mirrors.aibee.cn/tools/apache_hadoop/krb5.conf -O /etc/krb5.conf
COPY . /opt/ali/



