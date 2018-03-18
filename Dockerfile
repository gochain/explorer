FROM node:alpine

WORKDIR /explorer

RUN npm install -g grunt-cli
ADD package.* /explorer
RUN cd /explorer && npm install

ADD . /explorer

CMD ["npm", "start"]
