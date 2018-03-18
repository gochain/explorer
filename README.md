# GoChain Block Explorer

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
RPC_HOST - host of the gochain
RPC_PORT - port of the gochain
LISTEN - do not refill DB, just listen for new blocks
TERMINATE_DB - will terminate the block grabber once it gets to a block it has already stored in the DB

### Run:

`npm start`

Leave this running in the background to continuously fetch new blocks.

## Docker

Build:

`docker build -t gochain/explorer .`

Run:

 `docker run -e PORT=3000 -e RPC_HOST=$YOUR_RPC_HOST --net=host gochain/explorer`

* take in account that mongo should be running on same host (--net=host)
