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

MacOS: 

```
brew install mongodb
mongod --config /usr/local/etc/mongod.conf
```

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


# Block explorer API
## General endpoints
### Get stats
```
GET /api/stats
```

**Parameters:**
NONE

**Response:**
```javascript
{
    "updated_at":"2019-01-04T16:17:18.457Z","total_transactions_count":828310654,"last_week_transactions_count":19939999,"last_day_transactions_count":2856095
}
```

### Get the list of the most recent blocks
```
GET /api/blocks
```

**Parameters:**
limit - amount of items in the response
skip - number of items to skip


**Response:**
```javascript
{"blocks":[{"number":3265580,"created_at":"2019-01-04T16:23:39Z","miner":"0x6E2CAD5118b75420f7D73cfbD6db072523De8366","tx_count":0,"extra_data":"GoChain 5\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"},{"number":3265579,"created_at":"2019-01-04T16:23:34Z","miner":"0x7AeCEB5D345a01F8014a4320aB1F3D467c0C086a","tx_count":25,"extra_data":"GoChain 1\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"}]}
```

### Get specific block details
```
GET /api/blocks/{block_number}
```

**Parameters:**
block_number

**Response:**
```javascript
{"number":3265579,"gas_limit":136500000,"hash":"0x5c08a62b590c597713dcc8c68a1647ba086330c8bdf208543533367378622f38","created_at":"2019-01-04T16:23:34Z","parent_hash":"0x215f85d89767a89428649f64bacf135cf7491c96d0d8278fd7b912a900eaf914","tx_hash":"0x61a1c5d1609a3edb582fbbc32b64298ae5010fe99da039f4d7f5ff680559bcb8","gas_used":"525000","nonce":"0","miner":"0x7AeCEB5D345a01F8014a4320aB1F3D467c0C086a","tx_count":25,"difficulty":5,"total_difficulty":0,"sha3_uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","extra_data":"GoChain 1\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"}
```

### Get the specific block transactons
```
GET /api/blocks/{block_number}/transactions
```

**Parameters:**
block_number
limit - amount of items in the response
skip - number of items to skip


**Response:**
```javascript
{"transactions":[{"tx_hash":"0x9e5580ba859c11a4be5f62907d54063b23ccf5a9ca996b4d2a49f34005ff06dd","to":"0x44D63da717F5Cb2f74B3CFa9e02d633479CB1100","from":"0xB93901B9413a08DA4E90a2264d12B3eadB8dCA82","contract_address":"","value":"13","gas_price":"3970867867","gas_fee":"988825516240340","gas_limit":249020,"block_number":3265579,"nonce":"180319","block_hash":"0x5c08a62b590c597713dcc8c68a1647ba086330c8bdf208543533367378622f38","created_at":"2019-01-04T16:23:34Z","input_data":""}]}
```

### Get the specific address details
```
GET /api/address/{address_hash}
```

**Parameters:**
address_hash

**Response:**
```javascript
{"address":"0xeA224a724c08Ea862959AfE293C647c8671638f2","balance":"0.001059091847755743","balance_wei":"1059091847755743","updated_at":"2019-01-04T16:27:12.812Z","total_supply":"0","contract":false,"go20":false,"number_of_transactions":1800373}
```

### Get the specific address transactions
```
GET /api/address/{address_hash}/transactions
```

**Parameters:**
address_hash
limit - amount of items in the response
skip - number of items to skip
input_data_empty - true/false
from_time/to_time - filter for the transaction list (by created_at time)
**Response:**
```javascript
{"transactions":[{"tx_hash":"0x9e5580ba859c11a4be5f62907d54063b23ccf5a9ca996b4d2a49f34005ff06dd","to":"0x44D63da717F5Cb2f74B3CFa9e02d633479CB1100","from":"0xB93901B9413a08DA4E90a2264d12B3eadB8dCA82","contract_address":"","value":"13","gas_price":"3970867867","gas_fee":"988825516240340","gas_limit":249020,"block_number":3265579,"nonce":"180319","block_hash":"0x5c08a62b590c597713dcc8c68a1647ba086330c8bdf208543533367378622f38","created_at":"2019-01-04T16:23:34Z","input_data":""}]}
```


### Get the specific contracts token holders (for go20 compatible contracts)
```
GET /api/address/{address_hash}/holders
```

**Parameters:**
address_hash
limit - amount of items in the response
skip - number of items to skip


**Response:**
```javascript
{"token_holders":[{"contract_address":"0x7bb44320d4af20D259Ef6DD259bF00E246ed052c","token_holder_address":"0x744062E05485c94e92F5FeDB8E4067fcA01709e1","balance":"1000000000000000000000","balance_int":1000,"updated_at":"2018-12-16T15:04:58.088Z"},{"contract_address":"0x7bb44320d4af20D259Ef6DD259bF00E246ed052c","token_holder_address":"0x9E049291f70B917cd332420d29afDea9F4d76696","balance":"0","balance_int":0,"updated_at":"2018-12-16T15:04:58.257Z"},{"contract_address":"0x7bb44320d4af20D259Ef6DD259bF00E246ed052c","token_holder_address":"0x0000000000000000000000000000000000000000","balance":"0","balance_int":0,"updated_at":"2018-12-12T17:49:10.073Z"}]}
```

### Get the specific contracts internal transactions (for go20 compatible contracts)
```
GET /api/address/{address_hash}/internal_transactions
```

**Parameters:**
address_hash
limit - amount of items in the response
skip - number of items to skip


**Response:**
```javascript
{"internal_transactions":[{"contract_address":"0x7bb44320d4af20D259Ef6DD259bF00E246ed052c","from_address":"0x9E049291f70B917cd332420d29afDea9F4d76696","to_address":"0x744062E05485c94e92F5FeDB8E4067fcA01709e1","value":"1000000000000000000000","block_number":2420510,"transaction_hash":"0xcc543bab87ce87be04c873cc38c4f7cf627aaa33d9fe8934be155ea89a5d07e2","updated_at":"2018-12-16T15:04:57.918Z","created_at":"0001-01-01T00:00:00Z"},{"contract_address":"0x7bb44320d4af20D259Ef6DD259bF00E246ed052c","from_address":"0x0000000000000000000000000000000000000000","to_address":"0x9E049291f70B917cd332420d29afDea9F4d76696","value":"1000000000000000000000","block_number":2420503,"transaction_hash":"0x551b40fa88303e7b9e013da1b2d83a1a6d6f7677cbade8ae8027293004279841","updated_at":"2018-12-12T17:49:09.107Z","created_at":"0001-01-01T00:00:00Z"}]}
```

### Get the specific contracts details (where available)
```
GET /api/address/{address_hash}/contract
```

**Parameters:**
address_hash

**Response:**
```javascript
{}
```

### Get the specific transaction details
```
GET /api/transaction/{tx_hash}
```

**Parameters:**
tx_hash

**Response:**
```javascript
{"tx_hash":"0x0533322f4b48030bbc38498cfc4847b0bca7511454a351ce1bc700ad8eeb0ad7","to":"0xeA224a724c08Ea862959AfE293C647c8671638f2","from":"0x47ccA4785443B68a7363bF8b7995F95C46eb45a7","contract_address":"0x0000000000000000000000000000000000000000","value":"18","gas_price":"2125437567","gas_fee":"44634188907000","gas_limit":257573,"block_number":3265557,"nonce":"153174","block_hash":"0x565359100c9fd6418e6bf6395714728836ec9fbc0e3819f067cfe2a41059d01c","created_at":"2019-01-04T16:21:44Z","input_data":""}
```

### Get the richlist
```
GET /api/richlist
```

**Parameters:**
limit - amount of items in the response
skip - number of items to skip

**Response:**
```javascript
{"total_supply":"1022860075.000000000000000000","circulating_supply":"63410232.031792708000000000","rankings":[{"address":"0x78cb510135787f42a23aD46998eB16B756111559","balance":"19983424.526977560219899047","balance_wei":"19983424526977560219899047","updated_at":"2019-01-04T16:28:28.72Z","total_supply":"0","contract":false,"go20":false,"number_of_transactions":0}]}
```
