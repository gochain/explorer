# GoChain Block Explorer
[![CircleCI](https://circleci.com/gh/gochain-io/explorer.svg?style=svg)](https://circleci.com/gh/gochain-io/explorer)

Simple interface for exploring the GoChain blockchain.

## Local installation

Clone the repo

`git clone https://github.com/gochain-io/explorer`

Download [Nodejs and npm](https://docs.npmjs.com/getting-started/installing-node "Nodejs install") if you don't have them

Install mongodb:

MacOS: 

```
brew install mongodb
mongod --config /usr/local/etc/mongod.conf
```

Ubuntu: `sudo apt-get install -y mongodb-org`

## Build

To run a local environment, you'll need to build the internal toolsets `grabber` and `server` whose binary files you will run as below.

To create these binaries and install dependencies, use the Makefile and view it for internals and other options:

`make build`

## Running 

1) seed local Mongo database with `grabber` (./grabber/grabber)
```sh
> grabber help
NAME:
   grabber - Grabber populates a mongo database with explorer data.

USAGE:
   grabber [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rpc-url value, -u value         rpc api url (default: "https://rpc.gochain.io")
   --mongo-url value, -m value       mongo connection url (default: "127.0.0.1:27017")
   --mongo-dbname value, --db value  mongo database name (default: "blocks")
   --log value, -l value             loglevel debug/info/warn/fatal (default: "info")
   --start-from value, -s value      refill from this block (default: 0)
   --help, -h                        show help
   --version, -v                     print the version
```

2) run `server` (./server/server) (point to the same database name that you selected for seeding)
```sh
> server help
NAME:
   server - Server serves the explorer web interface, backed by a mongo database.

USAGE:
   server [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rpc-url value, -u value         rpc api url (default: "https://rpc.gochain.io")
   --mongo-url value, -m value       mongo connection url (default: "127.0.0.1:27017")
   --mongo-dbname value, --db value  mongo database name (default: "blocks")
   --log value, -l value             loglevel debug/info/warn/fatal, default is Info (default: "info")
   --dist value, -d value            folder that should be served (default: "../dist/explorer/")
   --recaptcha value, -r value       secret key for google recaptcha v3
   --help, -h                        show help
   --version, -v                     print the version
```

3) launch the web application

`cd frontend && npm start`


_At this point a local version of the application should be available on localhost port 4200. Optionally utilize Docker as below._


## Docker

Build:

`make docker`

Run:

```sh
docker run --net=host gochain/explorer grabber [flags]
docker run --net=host gochain/explorer server [flags]
```

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
```json
{  
   "updated_at":"2019-01-04T16:17:18.457Z",
   "total_transactions_count":828310654,
   "last_week_transactions_count":19939999,
   "last_day_transactions_count":2856095
}
```

### Get the list of the most recent blocks
```
GET /api/blocks
```

**Parameters:**

- limit - amount of items in the response
- skip - number of items to skip


**Response:**
```json
{
  "blocks": [
    {
      "number": 3265580,
      "created_at": "2019-01-04T16:23:39Z",
      "miner": "0x6E2CAD5118b75420f7D73cfbD6db072523De8366",
      "tx_count": 0,
      "extra_data": "GoChain 5\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"
    },
    {
      "number": 3265579,
      "created_at": "2019-01-04T16:23:34Z",
      "miner": "0x7AeCEB5D345a01F8014a4320aB1F3D467c0C086a",
      "tx_count": 25,
      "extra_data": "GoChain 1\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"
    }
  ]
}
```

### Get specific block details
```
GET /api/blocks/{block_number}
```

**Parameters:**

- block_number

**Response:**
```json
{
  "number": 3265579,
  "gas_limit": 136500000,
  "hash": "0x5c08a62b590c597713dcc8c68a1647ba086330c8bdf208543533367378622f38",
  "created_at": "2019-01-04T16:23:34Z",
  "parent_hash": "0x215f85d89767a89428649f64bacf135cf7491c96d0d8278fd7b912a900eaf914",
  "tx_hash": "0x61a1c5d1609a3edb582fbbc32b64298ae5010fe99da039f4d7f5ff680559bcb8",
  "gas_used": "525000",
  "nonce": "0",
  "miner": "0x7AeCEB5D345a01F8014a4320aB1F3D467c0C086a",
  "tx_count": 25,
  "difficulty": 5,
  "total_difficulty": 0,
  "sha3_uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
  "extra_data": "GoChain 1\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"
}
```

### Get the specific block transactons
```
GET /api/blocks/{block_number}/transactions
```

**Parameters:**

- block_number
- limit - amount of items in the response
- skip - number of items to skip


**Response:**
```json
{
  "transactions": [
    {
      "tx_hash": "0x9e5580ba859c11a4be5f62907d54063b23ccf5a9ca996b4d2a49f34005ff06dd",
      "to": "0x44D63da717F5Cb2f74B3CFa9e02d633479CB1100",
      "from": "0xB93901B9413a08DA4E90a2264d12B3eadB8dCA82",
      "contract_address": "",
      "value": "13",
      "gas_price": "3970867867",
      "gas_fee": "988825516240340",
      "gas_limit": 249020,
      "block_number": 3265579,
      "nonce": "180319",
      "block_hash": "0x5c08a62b590c597713dcc8c68a1647ba086330c8bdf208543533367378622f38",
      "created_at": "2019-01-04T16:23:34Z",
      "input_data": ""
    }
  ]
}
```

### Get the specific address details
```
GET /api/address/{address_hash}
```

**Parameters:**

- address_hash

**Response:**
```json
{
  "address": "0x1997eF6BeE5d61979E63d0c6b40F6d185Ab1156D",
  "balance": "0.000000000000000000",
  "balance_wei": "0",
  "updated_at": "2019-05-11T00:53:47.478+06:00",
  "token_name": "Example Fixed Supply Token",
  "token_symbol": "FIXED",
  "decimals": 18,
  "total_supply": "1000000000000000000000000",
  "contract": true,
  "erc_types": ["Go20","Go20Detailed"],
  "interfaces": null,
  "number_of_transactions": 3,
  "number_of_token_holders": 3,
  "number_of_internal_transactions": 2,
}
```

### Get the specific address transactions
```
GET /api/address/{address_hash}/transactions
```

**Parameters:**

- address_hash
- limit - amount of items in the response
- skip - number of items to skip
- from_time/to_time - filter for the transaction list (by created_at time)

**Response:**
```json
{
  "transactions": [
    {
      "tx_hash": "0x9e5580ba859c11a4be5f62907d54063b23ccf5a9ca996b4d2a49f34005ff06dd",
      "to": "0x44D63da717F5Cb2f74B3CFa9e02d633479CB1100",
      "from": "0xB93901B9413a08DA4E90a2264d12B3eadB8dCA82",
      "contract_address": "",
      "value": "13",
      "gas_price": "3970867867",
      "gas_fee": "988825516240340",
      "gas_limit": 249020,
      "block_number": 3265579,
      "nonce": "180319",
      "block_hash": "0x5c08a62b590c597713dcc8c68a1647ba086330c8bdf208543533367378622f38",
      "created_at": "2019-01-04T16:23:34Z",
      "input_data": ""
    }
  ]
}
```


### Get the specific contracts token holders (for go20 compatible contracts)
```
GET /api/address/{address_hash}/holders
```

**Parameters:**

- address_hash
- limit - amount of items in the response
- skip - number of items to skip


**Response:**
```json
{
  "token_holders": [
    {
      "contract_address": "0x7bb44320d4af20D259Ef6DD259bF00E246ed052c",
      "token_holder_address": "0x744062E05485c94e92F5FeDB8E4067fcA01709e1",
      "balance": "1000000000000000000000",
      "balance_int": 1000,
      "updated_at": "2018-12-16T15:04:58.088Z"
    },
    {
      "contract_address": "0x7bb44320d4af20D259Ef6DD259bF00E246ed052c",
      "token_holder_address": "0x9E049291f70B917cd332420d29afDea9F4d76696",
      "balance": "0",
      "balance_int": 0,
      "updated_at": "2018-12-16T15:04:58.257Z"
    },
    {
      "contract_address": "0x7bb44320d4af20D259Ef6DD259bF00E246ed052c",
      "token_holder_address": "0x0000000000000000000000000000000000000000",
      "balance": "0",
      "balance_int": 0,
      "updated_at": "2018-12-12T17:49:10.073Z"
    }
  ]
}
```

### Get the specific contracts internal transactions (for go20 compatible contracts)
```
GET /api/address/{address_hash}/internal_transactions
```

**Parameters:**

- address_hash
- limit - amount of items in the response
- skip - number of items to skip
- from_address - the transaction source address
- to_address - the transaction destination address


**Response:**
```json
{
  "internal_transactions": [
    {
      "contract_address": "0x7bb44320d4af20D259Ef6DD259bF00E246ed052c",
      "from_address": "0x9E049291f70B917cd332420d29afDea9F4d76696",
      "to_address": "0x744062E05485c94e92F5FeDB8E4067fcA01709e1",
      "value": "1000000000000000000000",
      "block_number": 2420510,
      "transaction_hash": "0xcc543bab87ce87be04c873cc38c4f7cf627aaa33d9fe8934be155ea89a5d07e2",
      "updated_at": "2018-12-16T15:04:57.918Z",
      "created_at": "0001-01-01T00:00:00Z"
    },
    {
      "contract_address": "0x7bb44320d4af20D259Ef6DD259bF00E246ed052c",
      "from_address": "0x0000000000000000000000000000000000000000",
      "to_address": "0x9E049291f70B917cd332420d29afDea9F4d76696",
      "value": "1000000000000000000000",
      "block_number": 2420503,
      "transaction_hash": "0x551b40fa88303e7b9e013da1b2d83a1a6d6f7677cbade8ae8027293004279841",
      "updated_at": "2018-12-12T17:49:09.107Z",
      "created_at": "0001-01-01T00:00:00Z"
    }
  ]
}
```

### Get the specific contracts details (where available)
```
GET /api/address/{address_hash}/contract
```

**Parameters:**

- address_hash

**Response:**
```json
{}
```

### Get the specific transaction details
```
GET /api/transaction/{tx_hash}
```

**Parameters:**

- tx_hash

**Response:**
```json
{
  "tx_hash": "0x0533322f4b48030bbc38498cfc4847b0bca7511454a351ce1bc700ad8eeb0ad7",
  "to": "0xeA224a724c08Ea862959AfE293C647c8671638f2",
  "from": "0x47ccA4785443B68a7363bF8b7995F95C46eb45a7",
  "contract_address": "0x0000000000000000000000000000000000000000",
  "value": "18",
  "gas_price": "2125437567",
  "gas_fee": "44634188907000",
  "gas_limit": 257573,
  "block_number": 3265557,
  "nonce": "153174",
  "block_hash": "0x565359100c9fd6418e6bf6395714728836ec9fbc0e3819f067cfe2a41059d01c",
  "created_at": "2019-01-04T16:21:44Z",
  "input_data": ""
}
```

### Get the richlist
```
GET /api/richlist
```

**Parameters:**

- limit - amount of items in the response
- skip - number of items to skip

**Response:**
```json
{
  "total_supply": "1022860075.000000000000000000",
  "circulating_supply": "63410232.031792708000000000",
  "rankings": [
    {
      "address": "0x78cb510135787f42a23aD46998eB16B756111559",
      "balance": "19983424.526977560219899047",
      "balance_wei": "19983424526977560219899047",
      "updated_at": "2019-02-04T16:28:28.72Z",
      "total_supply": "0",
      "contract": false,
      "number_of_transactions": 0
    }
  ]
}
```
