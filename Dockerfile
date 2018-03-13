FROM node

WORKDIR /explorer

ADD . /explorer
RUN npm install

CMD ["/bin/sh", "/explorer/start.sh"]

EXPOSE 8000
