require('../db.js');
var etherUnits = require("../lib/etherUnits.js");
var BigNumber = require('bignumber.js');

var fs = require('fs');

var Web3 = require('web3');

var mongoose = require('mongoose');
var Block = mongoose.model('Block');
var Transaction = mongoose.model('Transaction');

var grabBlocks = function () {
    var web3 = new Web3(new Web3.providers.HttpProvider('http://' + process.env.RPC_HOST.toString() + ':' +
        process.env.RPC_PORT.toString()));
    listenBlocks(web3);
    grabBlock(web3, { 'start': 0, 'end': 'latest' }, false);
}

var listenBlocks = function (web3) {
    console.log("Started listening for a new blocks")
    var newBlocks = web3.eth.filter("latest");
    newBlocks.watch(function (error, log) {
        if (error) {
            console.log('Error: ' + error);
            newBlocks.stopWatching();
            grabBlocks();
            console.log('Stopped watching, restarting filter');
        } else if (log == null) {
            console.log('Warning: null block hash');
        } else {
            console.log("Got new hash:", log);
            grabBlock(web3, log, true);
        }

    });
}

var grabBlock = function (web3, blockHashOrNumber, listening) {
    var desiredBlockHashOrNumber;
    // check if done
    if (blockHashOrNumber == undefined) {
        return;
    }

    if (typeof blockHashOrNumber === 'object') {
        if ('start' in blockHashOrNumber && 'end' in blockHashOrNumber) {
            desiredBlockHashOrNumber = blockHashOrNumber.end;
        }
        else {
            console.log('Error: Aborted becasue found a interval in blocks ' +
                'array that doesn\'t have both a start and end.');
            // process.exit(9);
        }
    }
    else {
        desiredBlockHashOrNumber = blockHashOrNumber;
    }

    if (web3.isConnected()) {
        web3.eth.getBlock(desiredBlockHashOrNumber, true, function (error, blockData) {
            if (error) {
                console.log('Warning: error on getting block with hash/number: ' +
                    desiredBlockHashOrNumber + ': ' + error);
            }
            else if (blockData == null) {
                console.log('Warning: null block data received from the block with hash/number: ' +
                    desiredBlockHashOrNumber);
            }
            else {
                checkBlockDBExistsThenWrite(blockData);
                if (listening == true)
                    return;
                if ('hash' in blockData && 'number' in blockData) {
                    // If currently working on an interval (typeof blockHashOrNumber === 'object') and
                    // the block number or block hash just grabbed isn't equal to the start yet:
                    // then grab the parent block number (<this block's number> - 1). Otherwise done
                    // with this interval object (or not currently working on an interval)
                    // -> so move onto the next thing in the blocks array.
                    if (typeof blockHashOrNumber === 'object' &&
                        (
                            (typeof blockHashOrNumber['start'] === 'string' && blockData['hash'] !== blockHashOrNumber['start']) ||
                            (typeof blockHashOrNumber['start'] === 'number' && blockData['number'] > blockHashOrNumber['start'])
                        )
                    ) {
                        blockHashOrNumber['end'] = blockData['number'] - 1;
                        checkBlockDBExistsThenGrab(web3, blockHashOrNumber);
                    }
                }
                else {
                    console.log('Error: No hash or number was found for block: ' + blockHashOrNumber);
                    // process.exit(9);
                }
            }
        }.bind({ listening: listening }));
    }
    else {
        console.log('Error: Aborted due to web3 is not connected when trying to ' +
            'get block ' + desiredBlockHashOrNumber);
        // process.exit(9);
    }
}


var writeBlockToDB = function (blockData) {
    blockData.transactionsCount = blockData.transactions.length
    return new Block(blockData).save(function (err, block, count) {
        if (typeof err !== 'undefined' && err) {
            if (err.code == 11000) {
                console.log('Skip: Duplicate key ' +
                    blockData.number.toString() + ': ' +
                    err);
            } else {
                console.log('Error: Aborted due to error on ' +
                    'block number ' + blockData.number.toString() + ': ' +
                    err);
                // process.exit(9);
            }
        } else {

            console.log('DB successfully written block number ' +
                blockData.number.toString());
        }
    });
}

var checkBlockDBExistsThenGrab = function (web3, blockHashOrNumber) {
    Block.find({ number: blockHashOrNumber['end'] }, function (err, b) {
        if (!b.length) {
            grabBlock(web3, blockHashOrNumber, false);
        } else {
            console.log("Block exist, trying next", blockHashOrNumber['end']);
            blockHashOrNumber['end'] = blockHashOrNumber['end'] - 1;
            checkBlockDBExistsThenGrab(web3, blockHashOrNumber);
        }

    })
}

/**
  * Checks if the a record exists for the block number then ->
  *     if record exists: abort
  *     if record DNE: write a file for the block
  */
var checkBlockDBExistsThenWrite = function (blockData) {
    Block.find({ number: blockData.number }, function (err, b) {
        if (!b.length) {
            writeBlockToDB(blockData);
            writeTransactionsToDB(blockData);
        } else {
            console.log("Block found in db: ", blockData.number.toString());
        }
    })
}

/**
    Break transactions out of blocks and write to DB
**/

var writeTransactionsToDB = function (blockData) {
    var bulkOps = [];
    if (blockData.transactions.length > 0) {
        for (d in blockData.transactions) {
            var txData = blockData.transactions[d];
            txData.timestamp = blockData.timestamp;
            txData.value = etherUnits.toEther(new BigNumber(txData.value), 'wei');
            bulkOps.push(txData);
        }
        Transaction.collection.insert(bulkOps, function (err, tx) {
            if (typeof err !== 'undefined' && err) {
                if (err.code == 11000) {
                    console.log('Skip: Duplicate key ' +
                        err);
                } else {
                    console.log('Error: Aborted due to error: ' +
                        err);
                    // process.exit(9);
                }
            } else
                console.log('DB successfully written block: ' + blockData.number.toString() + ', transactions: ' +
                    blockData.transactions.length.toString());
        });
    }
}

/*
  Patch Missing Blocks
*/
var patchBlocks = function (config) {
    var web3 = new Web3(new Web3.providers.HttpProvider('http://' + process.env.RPC_HOST.toString() + ':' +
        process.env.RPC_PORT.toString()));

    // number of blocks should equal difference in block numbers
    var firstBlock = 0;
    var lastBlock = web3.eth.blockNumber;
    blockIter(web3, firstBlock, lastBlock, config);
}

var blockIter = function (web3, firstBlock, lastBlock, config) {
    // if consecutive, deal with it
    if (lastBlock < firstBlock)
        return;
    if (lastBlock - firstBlock === 1) {
        [lastBlock, firstBlock].forEach(function (blockNumber) {
            Block.find({ number: blockNumber }, function (err, b) {
                if (!b.length)
                    grabBlock(web3, firstBlock);
            });
        });
    } else if (lastBlock === firstBlock) {
        Block.find({ number: firstBlock }, function (err, b) {
            if (!b.length)
                grabBlock(web3, firstBlock);
        });
    } else {

        Block.count({ number: { $gte: firstBlock, $lte: lastBlock } }, function (err, c) {
            var expectedBlocks = lastBlock - firstBlock + 1;
            if (c === 0) {
                grabBlock(web3, { 'start': firstBlock, 'end': lastBlock });
            } else if (expectedBlocks > c) {
                console.log("Missing: " + JSON.stringify(expectedBlocks - c));
                var midBlock = firstBlock + parseInt((lastBlock - firstBlock) / 2);
                blockIter(web3, firstBlock, midBlock, config);
                blockIter(web3, midBlock + 1, lastBlock, config);
            } else
                return;
        })
    }
}


/** On Startup **/
// geth --rpc --rpcaddr "localhost" --rpcport "8545"  --rpcapi "eth,net,web3"

var config = {};

// set the default geth port if it's not provided
if (!('gethPort' in config) || (typeof process.env.RPC_PORT) !== 'number') {
    process.env.RPC_PORT = 8545; // default
}


grabBlocks();
// patchBlocks(config);
