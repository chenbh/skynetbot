FROM golang:alpine as builder

WORKDIR /src
COPY go.* /src/
RUN go mod download

COPY . /src
RUN go build -o /bot ./cmd

FROM alpine:3
RUN apk add ffmpeg
COPY --from=builder /bot /bot

ENTRYPOINT /bot
