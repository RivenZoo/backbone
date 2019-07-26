# Version: 1.0
FROM golang:1.12 as builder

RUN mkdir -p /data/app
ADD ./ /data/app
RUN go get -u github.com/RivenZoo/backbone/logger && cd /data/app && go build .

FROM golang:1.12

RUN mkdir -p /data/app

WORKDIR /data/app

COPY --from=builder /data/app /usr/local/bin

