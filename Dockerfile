# Build GoChain in a stock Go builder container
FROM golang:1.10-alpine as backend_builder
RUN apk --no-cache add build-base git bzr mercurial gcc linux-headers g++ make
ENV D=/go/src/github.com/gochain-io/explorer
RUN go get -u github.com/golang/dep/cmd/dep
ADD . $D
RUN cd $D && make backend && mkdir -p /tmp/gochain && cp $D/server/server /tmp/gochain/ && cp $D/grabber/grabber /tmp/gochain/

FROM node:8-alpine  as frontend_builder
WORKDIR /explorer
RUN apk add --no-cache make
ADD . /explorer
RUN npm install -g @angular/cli@6.0.8
RUN make frontend

FROM alpine:latest
WORKDIR /explorer
RUN apk add --no-cache ca-certificates
COPY --from=backend_builder /tmp/gochain/* /usr/local/bin/
COPY --from=frontend_builder /explorer/dist/* /explorer/

EXPOSE 8080

CMD [ "server","-d", "/explorer/" ]
