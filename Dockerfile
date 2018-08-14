# Build GoChain in a stock Go builder container
FROM golang:1.10-alpine as builder
RUN apk --no-cache add build-base git bzr mercurial gcc linux-headers
ENV D=/go/src/github.com/gochain-io/explorer
# RUN go get -u github.com/golang/dep/cmd/dep
# ADD Gopkg.* $D/
ADD . $D
RUN cd $D && make backend && mkdir -p /tmp/gochain && cp $D/server/server /tmp/gochain/ && cp $D/grabber/grabber /tmp/gochain/

# Pull all binaries into a second stage deploy alpine container
FROM node:alpine

WORKDIR /explorer

RUN apk add --no-cache ca-certificates
COPY --from=builder /tmp/gochain/* /usr/local/bin/

RUN npm install -g grunt-cli
RUN make buildfront
ADD dist /explorer