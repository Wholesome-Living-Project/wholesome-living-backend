FROM ubuntu:latest

# System dependencies
RUN apt-get update -y
RUN apt-get install -y wget build-essential

# Installation of golang
RUN wget -q https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

# Installation of air
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download
ENV PATH=$PATH:/root/go/bin

COPY . /app
WORKDIR app

EXPOSE 8080


CMD air
