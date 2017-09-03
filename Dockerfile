FROM golang:1.9-alpine AS build

RUN apk add --update curl git
RUN curl https://glide.sh/get | sh

RUN mkdir -p /go/src/github.com/gquental/pokedex
WORKDIR /go/src/github.com/gquental/pokedex

COPY glide.lock .
COPY glide.yaml .

RUN glide install

COPY . .

RUN mkdir /app

RUN go build -o /app/pokedex
RUN go build -o /app/pokedex-cli cli/main.go

CMD ["/app/pokedex"]
