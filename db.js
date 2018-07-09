var mongoose = require('mongoose');
var Schema = mongoose.Schema;

var Block = new Schema(
    {
        "number": { type: Number, index: { unique: true } },
        "hash": String,
        "parentHash": String,
        "nonce": String,
        "sha3Uncles": String,
        "logsBloom": String,
        "transactionsRoot": String,
        "stateRoot": String,
        "receiptRoot": String,
        "miner": { type: String, index: true },
        "difficulty": String,
        "totalDifficulty": String,
        "size": Number,
        "extraData": String,
        "gasLimit": Number,
        "gasUsed": Number,
        "timestamp": Number,
        "transactionsCount": Number,
        "uncles": [String]
    });

var Contract = new Schema(
    {
        "address": { type: String, index: { unique: true } },
        "creationTransaction": String,
        "contractName": String,
        "compilerVersion": String,
        "optimization": Boolean,
        "sourceCode": String,
        "abi": String,
        "byteCode": String
    }, { collection: "Contract" });

var Transaction = new Schema(
    {
        "hash": { type: String, index: true },
        "nonce": Number,
        "blockHash": String,
        "blockNumber": { type: Number, index: true },
        "transactionIndex": Number,
        "from": { type: String, index: true },
        "to": { type: String, index: true },
        "value": String,
        "gas": Number,
        "gasPrice": String,
        "timestamp": { type: Number, index: true },
        "input": String
    });

var Address = new Schema(
    {
        "address": { type: String, index: { unique: true } },
        "owner": String,
        "balance": String,
        "lastUpdated": { type: Number, index: true },
        "balanceDecimal": { type: Number, index: true },
        "transactionsCount": Number
    });

Transaction.index({ hash: 1, blockNumber: 1 }, { unique: true });
mongoose.model('Block', Block);
mongoose.model('Contract', Contract);
mongoose.model('Transaction', Transaction);
mongoose.model('Address', Address);
module.exports.Block = mongoose.model('Block');
module.exports.Contract = mongoose.model('Contract');
module.exports.Transaction = mongoose.model('Transaction');
module.exports.Address = mongoose.model('Address');

const serverOptions = {
    poolsize: 100,
    socketOptions: {
        keepAlive: 1,
        auto_reconnect: true,
        connectTimeoutMS: 6000000,
        socketTimeoutMS: 6000000
    }
};
mongoose.connect(process.env.MONGO_URI || 'mongodb://localhost/blockDB', {
    server: serverOptions
});
mongoose.set('debug', false);
