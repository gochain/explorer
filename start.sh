#!/usr/bin/env sh

: ${GETH_HOSTNAME:="localhost"}
: ${GETH_RPCPORT:=8545}

echo $GETH_HOSTNAME
sed -i -e "s/var GETH_HOSTNAME.*/var GETH_HOSTNAME = \"${GETH_HOSTNAME}\";/g" -e "s/var GETH_RPCPORT.*/var GETH_RPCPORT = ${GETH_RPCPORT}/g" app/app.js
head app/app.js
npm start
