# Developer
FROM golang:1.19-alpine as dev
RUN apk update && apk upgrade && \
    apk --update add gcc git make curl build-base && \
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s && \
    air -v && \
    curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s && \
    mkdir -p /usr/github.com/isaias-dgr/currency-tracker && \
    mkdir -p /usr/github.com/isaias-dgr/currency-tracker/tmp && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/pressly/goose/v3/cmd/goose@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest

WORKDIR /usr/github.com/isaias-dgr/currency-tracker
COPY . /usr/github.com/isaias-dgr/currency-tracker

# Builder
FROM golang:1.19-alpine as builder
RUN apk update && apk upgrade && \
    apk --update add gcc git make curl

RUN mkdir -p /usr/github.com/isaias-dgr/currency-tracker 

WORKDIR /usr/github.com/isaias-dgr/currency-tracker
COPY . /usr/github.com/isaias-dgr/currency-tracker
RUN go mod download && \ 
    go build -o /tmp/app/main src/app/main.go

# Distribution 
FROM alpine:latest as release

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir -p /src/app 
COPY migrate /migrate
WORKDIR /src/app 

EXPOSE 8080
COPY --from=builder /tmp/app/main /src/app
CMD /src/app/main