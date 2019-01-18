# Build GoChain in a stock Go builder container
FROM golang:alpine as backend_builder
RUN apk --no-cache add build-base git bzr mercurial gcc linux-headers g++ make
ENV D=/go/src/github.com/gochain-io/explorer
RUN go get -u github.com/golang/dep/cmd/dep
ADD . $D
RUN cd $D && make backend && mkdir -p /tmp/gochain && cp $D/server/server /tmp/gochain/ && cp $D/grabber/grabber /tmp/gochain/

FROM node:8-alpine  as frontend_builder
WORKDIR /explorer
RUN apk add --no-cache make git gcc g++ python
ADD . /explorer
RUN npm install -g @angular/cli@7.2.1
RUN make frontend

FROM ethereum/solc:stable as solc

FROM alpine:latest
WORKDIR /explorer
RUN apk add --no-cache ca-certificates
COPY --from=backend_builder /tmp/gochain/* /usr/local/bin/
COPY --from=solc /usr/bin/solc /usr/local/bin/
COPY --from=frontend_builder /explorer/dist/* /explorer/

EXPOSE 8080

CMD [ "server","-d", "/explorer/" ]
