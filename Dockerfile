# Build GoChain in a stock Go builder container
FROM golang:1-alpine as backend_builder
RUN apk --no-cache add build-base git mercurial gcc linux-headers g++ make
ENV D=/explorer
WORKDIR $D
# cache dependencies
ADD go.mod $D
ADD go.sum $D
RUN go mod download
ADD . $D
# build
RUN cd $D && make backend && mkdir -p /tmp/gochain && cp $D/server/server /tmp/gochain/ && cp $D/grabber/grabber /tmp/gochain/ && cp $D/admin/admin /tmp/gochain/

FROM node:10-alpine as frontend_builder
WORKDIR /explorer
RUN apk add --no-cache make git gcc g++ python
ADD . /explorer
RUN npm install -g @angular/cli
RUN make frontend

FROM alpine:latest
WORKDIR /explorer
RUN apk add --no-cache ca-certificates docker
COPY --from=backend_builder /tmp/gochain/* /usr/local/bin/
COPY --from=frontend_builder /explorer/front/dist/* /explorer/

EXPOSE 8080

CMD [ "server","-d", "/explorer/" ]
