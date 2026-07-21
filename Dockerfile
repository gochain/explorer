# Build GoChain in a stock Go builder container
FROM golang:1 as backend_builder
ENV D=/explorer
WORKDIR $D
# cache dependencies
ADD go.mod $D
ADD go.sum $D
RUN go mod download
ADD . $D
# build
RUN cd $D && make backend && mkdir -p /tmp/gochain && cp $D/server/server /tmp/gochain/ && cp $D/grabber/grabber /tmp/gochain/ && cp $D/admin/admin /tmp/gochain/

FROM node:26 as frontend_builder
ENV NODE_OPTIONS=--openssl-legacy-provider
WORKDIR /explorer
RUN apt-get update && apt-get install -y make git gcc g++ python3 && rm -rf /var/lib/apt/lists/*
ADD . /explorer
RUN npm install -g @angular/cli@latest
RUN make frontend

FROM ubuntu:latest
WORKDIR /explorer
RUN apt-get update && apt-get install -y ca-certificates docker.io && rm -rf /var/lib/apt/lists/*
COPY --from=backend_builder /tmp/gochain/* /usr/local/bin/
COPY --from=frontend_builder /explorer/front/dist/* /explorer/

EXPOSE 8080

CMD [ "server","-d", "/explorer/" ]
