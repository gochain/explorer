require('../db.js');
var etherUnits = require("../lib/etherUnits.js");
var BigNumber = require('bignumber.js');

var fs = require('fs');

var Web3 = require('web3');

var mongoose = require('mongoose');
var Block = mongoose.model('Block');
var Transaction = mongoose.model('Transaction');

var grabBlocks = function (web3) {
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
                checkBlockDBExistsThenWrite(web3, blockData);
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


var writeBlockToDB = function (web3, blockData) {
    blockData.transactionsCount = blockData.transactions.length
    return new Block(blockData).save(function (err, block, count) {
        if (typeof err !== 'undefined' && err) {
            if (err.code == 11000) {
                console.log('Skip: Duplicate key ' +
                    blockData.number.toString() + ': ' +
                    err);
                cleanupBlockAndTransactionsThenGrab(web3, blockData)
            } else {
                console.log('Error: Aborted due to error on ' +
                    'block number ' + blockData.number.toString() + ': ' +
                    err);
                // process.exit(9);
            }
        } else {

            console.log('DB successfully written block number ' +
                blockData.number.toString(), " block hash :" + blockData.hash);
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
  *     cleanup records for specific block and transactions
  */
var cleanupBlockAndTransactionsThenGrab = function (web3, blockData) {
    Block.remove({ number: blockData.number }, function (err) {
        if (typeof err !== 'undefined' && err) {
            console.log("Cannot remove block :", blockData.number, " Err:", err);
        } else {
            Transaction.remove({ blockNumber: blockData.number }, function (err) {
                if (typeof err !== 'undefined' && err) {
                    console.log("Cannot remove transactions for block :", blockData.number, " Err:", err);
                } else {
                    console.log("Transactions and block has been removed calling grab again:", blockData.number);
                    setTimeout(function () {
                        grabBlock(web3, blockData.number, false);
                    }, 3000);
                }
            });
        }
    });
}


/**
  * Checks if the a record exists for the block number then ->
  *     if record exists: abort
  *     if record DNE: write a file for the block
  */
var checkBlockDBExistsThenWrite = function (web3, blockData) {
    Block.find({ number: blockData.number }, function (err, b) {
        if (!b.length) {
            writeBlockToDB(web3, blockData);
            writeTransactionsToDB(blockData);
        } else {
            console.log("Block found in db: ", blockData.number.toString(), " block hash :" + blockData.hash);
        }
    })
}

/**
    Break transactions out of blocks and write to DB
**/

var writeTransactionsToDB = function (blockData) {

    if (blockData.transactions.length > 0) {
        console.log("Block: ", blockData.number.toString(), " trying to add transactions to db:", blockData.transactions.length, " block hash :" + blockData.hash);
        var chunkSize = 1000 //1000 transactions per block
        for (var i = 0; i < blockData.transactions.length; i += chunkSize) {
            chunk = blockData.transactions.slice(i, i + chunkSize);
            var bulkOps = [];
            for (d in chunk) {
                var txData = chunk[d];
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
                }
            });
        }
    }
}

// /*
//   Patch Missing Blocks
// */
// var patchBlocks = function (config) {
//     var web3 = new Web3(new Web3.providers.HttpProvider('http://' + process.env.RPC_HOST.toString() + ':' +
//         process.env.RPC_PORT.toString()));

//     // number of blocks should equal difference in block numbers
//     var firstBlock = 0;
//     var lastBlock = web3.eth.blockNumber;
//     blockIter(web3, firstBlock, lastBlock, config);
// }

// var blockIter = function (web3, firstBlock, lastBlock, config) {
//     // if consecutive, deal with it
//     if (lastBlock < firstBlock)
//         return;
//     if (lastBlock - firstBlock === 1) {
//         [lastBlock, firstBlock].forEach(function (blockNumber) {
//             Block.find({ number: blockNumber }, function (err, b) {
//                 if (!b.length)
//                     grabBlock(web3, firstBlock);
//             });
//         });
//     } else if (lastBlock === firstBlock) {
//         Block.find({ number: firstBlock }, function (err, b) {
//             if (!b.length)
//                 grabBlock(web3, firstBlock);
//         });
//     } else {

//         Block.count({ number: { $gte: firstBlock, $lte: lastBlock } }, function (err, c) {
//             var expectedBlocks = lastBlock - firstBlock + 1;
//             if (c === 0) {
//                 grabBlock(web3, { 'start': firstBlock, 'end': lastBlock });
//             } else if (expectedBlocks > c) {
//                 console.log("Missing: " + JSON.stringify(expectedBlocks - c));
//                 var midBlock = firstBlock + parseInt((lastBlock - firstBlock) / 2);
//                 blockIter(web3, firstBlock, midBlock, config);
//                 blockIter(web3, midBlock + 1, lastBlock, config);
//             } else
//                 return;
//         })
//     }
// }


// set the default geth port if it's not provided
var url = process.env.RPC_URL || 'http://localhost:8545'
console.log("CONNECTING TO:", url);
web3 = new Web3(new Web3.providers.HttpProvider(url));

grabBlocks(web3);
