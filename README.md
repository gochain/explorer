# GoChain Block Explorer
[![CircleCI](https://circleci.com/gh/gochain-io/explorer.svg?style=svg)](https://circleci.com/gh/gochain-io/explorer)

Simple interface for exploring the GoChain blockchain.

## Local installation

Clone the repo

`git clone https://github.com/gochain-io/explorer`

Download [Nodejs and npm](https://docs.npmjs.com/getting-started/installing-node "Nodejs install") if you don't have them

Install dependencies:

`npm install`

Install mongodb:

MacOS: `brew install mongodb`

Ubuntu: `sudo apt-get install -y mongodb-org`

## Populate the DB

This will fetch and parse the entire blockchain.

Basic settings(environment variables):

PORT - server port where UI is running
MONGO_URI - url to mongo DB ('mongodb://localhost/blockDB')
RPC_URL - URL of the gochain (default it http://localhost:8545)

### Run:

Modify `.env.sh` to include proper environment variables:

    export PORT=8000
    export MONGO_URI=mongodb://localhost:27017/blockDB
    export RPC_URL=https://testnet-rpc.gochain.io:443
    export RPC_IP=testnet-rpc.gochain.io

Then run `npm start`

Leave this running in the background to continuously fetch new blocks.

## Docker

Build:

`docker build -t gochain/explorer .`

Run:

 `docker run -e PORT=3000 -e RPC_HOST=$YOUR_RPC_HOST --net=host gochain/explorer`

* take in account that mongo should be running on same host (--net=host)
