require('../db.js');
var etherUnits = require("../lib/etherUnits.js");
var BigNumber = require('bignumber.js');

var fs = require('fs');

var Web3 = require('web3');

var mongoose = require('mongoose');
var Block = mongoose.model('Block');
var Transaction = mongoose.model('Transaction');
var Address = mongoose.model('Address');

var parseArgs = function () {
    return process.argv.slice(2).reduce((acc, arg) => {
        let [k, v = true] = arg.split('=')
        acc[k] = v
        return acc
    }, {})
}
var grabBlocks = function (web3) {
    args = parseArgs();
    var argsNotSet = Object.keys(args).length == 0
    console.log("in grab blocks");
    console.log("argsNotSet:", argsNotSet, " args:", args);

    if (args["listen"] || argsNotSet)
        listenBlocks(web3);
    if (args["refill"] || argsNotSet)
        grabBlock(web3, { 'start': 0, 'end': 'latest' }, false);
    if (args["richlist"] || argsNotSet)
        updateAddressesBalance(web3);
}

var updateAddressesBalance = function (web3, latestUpdate) {
    if (!latestUpdate) { latestUpdate = 0 }
    console.log("updateStartedAt", latestUpdate);
    var updateStartedAt = Date.now() / 1000 //need to divide by 1000 because transactions and blocks are using this format
    var genesisAllocAddress = []
    try {
        res = web3.currentProvider.send({ jsonrpc: "2.0", method: "eth_genesisAlloc", id: new Date().getTime() })
        genesisAllocAddress = Object.keys(res["result"])
    } catch (e) {
        console.log("Cannot get genesis");
    }
    Transaction.distinct("to", { timestamp: { $gte: latestUpdate } }, function (err, toAdresses) {
        if (!err) {
            Transaction.distinct("from", { timestamp: { $gte: latestUpdate } }, function (err, fromAdresses) {
                if (!err) {
                    Block.distinct("miner", { timestamp: { $gte: latestUpdate } }, function (err, miners) {
                        if (!err) {
                            var adressesToUpdate = toAdresses.concat(fromAdresses).concat(miners).concat(genesisAllocAddress);
                            uniqueArray = adressesToUpdate.filter(function (elem, pos) {
                                return adressesToUpdate.indexOf(elem) == pos && elem && genesisAllocAddress.indexOf(elem) == -1;
                            });
                            console.log("Got list of addresses to update:", uniqueArray.length);
                            var i = 0;
                            var len = uniqueArray.length;
                            function iter() {
                                if (i < len) {
                                    // console.log("Checking i:", i, " address:", uniqueArray[i]);
                                    updateAddressBalance(uniqueArray[i], web3, updateStartedAt);
                                    i++;
                                    setImmediate(iter);
                                } else {
                                    setTimeout(function () {
                                        updateAddressesBalance(web3, updateStartedAt);
                                    }, 300000);//fire after 5 minutes of run
                                }
                            }
                            iter();
                        } else {
                            console.log("Cannot make distinct for blocks:", err)
                        }
                    })
                } else {
                    console.log("Cannot make distinct for the from field of transactions:", err)
                }
            })
        } else {
            console.log("Cannot make distinct for the to field of transactions:", err)
        }
    })
}

var updateAddressBalance = function (address, web3, updateStartedAt) {
    try {
        var balance = web3.fromWei(web3.eth.getBalance(address));
        console.log("Got balance for address:", address, " balance:", balance.toNumber());
        Address.findOneAndUpdate({ address: address }, { $set: { balance: balance.toString(), balanceDecimal: balance.toNumber(), lastUpdated: updateStartedAt } }, { upsert: true }, function (err, doc) {
            if (err) {
                console.log("Something wrong when updating address:", address, err);
            }
            // console.log("Balance for address updated:", address);
        });
    } catch (error) {
        console.log("Exception while checking the address, retrying in 10 seconds:", address, " exception:", error);
        setTimeout(function () {
            console.log("Grabbing address balance after sleep :", address);
            updateAddressBalance(address, web3, updateStartedAt);
        }, Math.ceil(Math.random() * (60000 - 10000) + 10000));
    }
}

var listenBlocks = function (web3) {
    console.log("Started listening for a new blocks")
    var newBlocks
    try {
        newBlocks = web3.eth.filter("latest");
    } catch (e) {
        if (e instanceof TypeError) {
            grabBlocks(web3);
        } else {
            console.log("Exception:", e)
            throw (e);
        }
    }
    newBlocks.watch(function (error, hash) {
        if (error) {
            console.log('Error: ' + error);
            newBlocks.stopWatching();
            grabBlocks(web3);
            console.log('Stopped watching, restarting filter');
        } else if (hash == null) {
            console.log('Warning: null block hash');
        } else {
            console.log("Got new hash:", hash);
            web3.eth.getBlock(hash, false, function (error, block) {
                if (error) {
                    console.log('Warning: error on getting block with hash/number: ' +
                        hash + ': ' + error);
                } else {
                    setTimeout(function () {
                        grabBlock(web3, block.number, true);
                    }, 5000);
                }
            });

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
        }
    }
    else {
        desiredBlockHashOrNumber = blockHashOrNumber;
    }

    if (web3.isConnected()) {
        web3.eth.getBlock(desiredBlockHashOrNumber, true, function (error, blockData) {
            if (error) {
                console.log("Waiting X seconds and trying to grab again block id:", desiredBlockHashOrNumber, " error:", error);
                setTimeout(function () {
                    console.log("Grabbing after sleep block:", desiredBlockHashOrNumber, " error:", error);
                    grabBlock(web3, desiredBlockHashOrNumber, listening);
                }, Math.ceil(Math.random() * (60000 - 10000) + 10000));
            }
            else if (blockData == null) {
                console.log('Warning: null block data received from the block with hash/number: ' +
                    desiredBlockHashOrNumber);
                // process.exit(9);
            }
            else {
                checkBlockDBExistsThenWrite(web3, blockData);
                checkParentBlock(web3, blockData, listening);
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
                } else {
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

var checkParentBlock = function (web3, blockData, recursively) {
    if (blockData.number) {
        var parentBlockNumber = blockData.number - 1;
        Block.findOne({ number: parentBlockNumber }, function (err, b) {
            if (err) {
                console.log("Cannot find block in db in checkParentBlock:", err);
                grabBlock(web3, parentBlockNumber, recursively);
            } else {
                if (b) {
                    if (b && blockData.parentHash != b.hash) {
                        console.log("BLOCK IS BAD, GRABBING IT", parentBlockNumber);
                        cleanupBlockAndTransactionsThenGrab(web3, parentBlockNumber);
                        if (recursively) {
                            console.log("Checking recursively for:", parentBlockNumber);
                            Block.findOne({ number: parentBlockNumber - 1 }, function (err, b) {
                                if (err) {
                                    console.log("Cannot find the block in the db while checking recursively:", err)
                                } else {
                                    if (b) {
                                        checkParentBlock(web3, blockData, recursively);
                                    }
                                }
                            });
                        }

                    }
                } else {
                    // console.log("Cannot find parent block in db:", parentBlockNumber);
                    if (recursively)
                        grabBlock(web3, blockData.parentHash, recursively);
                }
            }
        });

    }
}
var writeBlockToDB = function (web3, blockData) {
    blockData.transactionsCount = blockData.transactions.length
    return new Block(blockData).save(function (err, block, count) {
        if (typeof err !== 'undefined' && err) {
            if (err.code == 11000) {
                console.log('Skip: Duplicate key in block ' +
                    blockData.number.toString());
                Block.findOne({ number: blockData.number.toString() }, function (err, b) {
                    if (b && blockData.hash.toString() != b.hash.toString()) {
                        console.log("HASH FROM API:", blockData.hash.toString(), " HASH IN DB:", b.hash.toString());
                        cleanupBlockAndTransactionsThenGrab(web3, blockData.number)
                    }
                })
                checkParentBlock(web3, blockData, false);
                // cleanupBlockAndTransactionsThenGrab(web3, blockData.number)
            } else {
                console.log('Error: Aborted due to error on ' +
                    'block number ' + blockData.number.toString() + ': ' +
                    err);
            }
        } else {

            console.log('DB successfully written block number ' +
                blockData.number.toString(), " block hash :" + blockData.hash);
        }
    });
}

var checkBlockDBExistsThenGrab = function (web3, blockHashOrNumber) {
    Block.findOne({ number: blockHashOrNumber['end'] }, function (err, b) {
        if (!b) {
            grabBlock(web3, blockHashOrNumber, false);
        } else {
            checkParentBlock(web3, b, false);
            checkTransactionsNumber(web3, b);
            if (b.number % 10000 == 0) { console.log("Block exist, trying next", blockHashOrNumber['end']) }
            blockHashOrNumber['end'] = blockHashOrNumber['end'] - 1;
            if (blockHashOrNumber['end'] > blockHashOrNumber['start']) {
                checkBlockDBExistsThenGrab(web3, { 'start': blockHashOrNumber['start'], 'end': blockHashOrNumber['end'] });
            } else {
                console.log("Reached end, starting from the beginning")
                grabBlock(web3, { 'start': blockHashOrNumber['start'], 'end': 'latest' }, false);
            }
        }

    })
}
var checkTransactionsNumber = function (web3, block) {
    Transaction.count({ blockNumber: block.number }, function (err, transactionsNumber) {
        if (err) {
            console.error(err);
        } else {
            if (transactionsNumber != block.transactionsCount) {
                console.log("Number of transactions from tx collection:", transactionsNumber, " from block:", block.transactionsCount, " block number:", block.number);
                cleanupBlockAndTransactionsThenGrab(web3, block.number)
            }
        }
    });
}

/**
  *     cleanup records for specific block and transactions
  */
var cleanupBlockAndTransactionsThenGrab = function (web3, blockNumber) {
    Block.remove({ number: blockNumber }, function (err) {
        if (typeof err !== 'undefined' && err) {
            console.log("Cannot remove block :", blockNumber, " Err:", err);
        } else {
            Transaction.remove({ blockNumber: blockNumber }, function (err) {
                if (typeof err !== 'undefined' && err) {
                    console.log("Cannot remove transactions for block :", blockNumber, " Err:", err);
                } else {
                    console.log("Transactions and block have been removed calling grab again:", blockNumber);
                    grabBlock(web3, blockNumber, true);
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
        console.log("Block: ", blockData.number.toString(), " trying to add transactions to db:", blockData.transactions.length, " block hash :" + blockData.hash, "Miner: " + blockData.miner.toString());
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
                        console.log('Skip: Duplicate key in transactions ' +
                            err);
                        cleanupBlockAndTransactionsThenGrab(web3, blockData.number);
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



// set the default geth port if it's not provided
var url = process.env.RPC_URL || 'http://localhost:8545'
console.log("CONNECTING TO:", url);
web3 = new Web3(new Web3.providers.HttpProvider(url));

grabBlocks(web3);
