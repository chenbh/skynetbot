FROM golang:alpine as builder

COPY go.* /src/

WORKDIR /src

RUN go mod download

COPY . /src

RUN go build -o /bot ./cmd

ENTRYPOINT [/bot]
