# Build powerplay in a stock Go builder container
FROM golang:1.10-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git


RUN git clone https://github.com/playmakerchain/powerplay.git
WORKDIR  /go/powerplay
RUN git checkout $(git describe --tags `git rev-list --tags --max-count=1`)
RUN make dep && make powerplay

# Pull powerplay into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/thor/bin/powerplay /usr/local/bin/

EXPOSE 2843 11235 11235/udp
ENTRYPOINT ["powerplay"]