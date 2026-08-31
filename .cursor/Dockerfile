FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        golang-go \
        postgresql \
        postgresql-client \
        python3 \
    && rm -rf /var/lib/apt/lists/*

# Go modules expect a recent toolchain; Ubuntu 24.04 ships Go 1.22.
ENV GOPATH=/home/ubuntu/go
ENV PATH=$PATH:/usr/lib/go-1.22/bin:$GOPATH/bin

USER ubuntu
