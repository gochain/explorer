#!/usr/bin/env node

/*
    Endpoint for client to talk to etc node
*/

var Web3 = require("web3");
var web3;

var BigNumber = require('bignumber.js');
var etherUnits = require(__lib + "etherUnits.js")

var getLatestBlocks = require('./index').getLatestBlocks;
var filterBlocks = require('./filters').filterBlocks;
var filterTrace = require('./filters').filterTrace;

var totalSupplyCached = 1000000000;
var circulatingSupplyCached = 500000000;


if (typeof web3 !== "undefined") {
  web3 = new Web3(web3.currentProvider);
} else {
  var url = process.env.RPC_URL || 'http://localhost:8545'
  web3 = new Web3(new Web3.providers.HttpProvider(url));
}

if (web3.isConnected())
  console.log("Web3 connection established");
else
  throw "No connection";


var newBlocks = web3.eth.filter("latest");
var newTxs = web3.eth.filter("pending");

exports.totalSupply = function (callback) {
  web3.currentProvider.sendAsync({
    jsonrpc: "2.0",
    method: "eth_totalSupply",
    params: ["latest"],
    id: new Date().getTime()
  }, function (error, result) {
    if (!error && result && result["result"]) {
      wei = new BigNumber(parseInt(result["result"], 16));
      go = etherUnits.toEther(wei, "wei");
      totalSupplyCached = go;
      callback(go);
    } else {
      console.log("Got error from API in totalSupply:" + error);
      callback(totalSupplyCached);
    }
  });
}

exports.circulatingSupply = function (callback) {
  this.totalSupply(function (total) {
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "eth_genesisAlloc",
      id: new Date().getTime()
    }, function (error, result) {
      if (!error && result && result["result"]) {
        allocGo = new BigNumber(Object.keys(result["result"]).map(function (v) {          
          return web3.fromWei(web3.eth.getBalance(v));
        }).reduce((a, b) => a + b, 0));        
        circulatingSupplyCached = total - allocGo;
        callback(total - allocGo);
      } else {
        console.log("Got error from API in circulatingSupply:" + error);
        callback(totalSupplyCached);
      }
    });
  });
}


exports.data = function (req, res) {
  console.log(req.body)

  if ("tx" in req.body) {
    var txHash = req.body.tx.toLowerCase();

    web3.eth.getTransaction(txHash, function (err, tx) {
      if (err || !tx) {
        console.error("TxWeb3 error :" + err)
        res.write(JSON.stringify({ "error": true }));
        res.end();
      } else {
        var ttx = tx;
        const gasPrice = new BigNumber(tx.gasPrice);
        const gas = new BigNumber(tx.gas);

        ttx.value = etherUnits.toEther(tx.value, "wei");
        ttx.actualGasCost = etherUnits.toEther(gas.multipliedBy(gasPrice), "wei");
        //get timestamp from block
        var block = web3.eth.getBlock(tx.blockNumber, function (err, block) {
          if (!err && block)
            ttx.timestamp = block.timestamp;
          ttx.isTrace = (ttx.input != "0x");
          res.write(JSON.stringify(ttx));
          res.end();
        });
      }
    });

  } else if ("tx_trace" in req.body) {
    var txHash = req.body.tx_trace.toLowerCase();

    web3.trace.transaction(txHash, function (err, tx) {
      if (err || !tx) {
        console.error("TraceWeb3 error :" + err)
        res.write(JSON.stringify({ "error": true }));
      } else {
        res.write(JSON.stringify(filterTrace(tx)));
      }
      res.end();
    });
  } else if ("addr_trace" in req.body) {
    var addr = req.body.addr_trace.toLowerCase();
    // need to filter both to and from
    // from block to end block, paging "toAddress":[addr],
    // start from creation block to speed things up
    // TODO: store creation block
    var filter = { "fromBlock": "0x1d4c00", "toAddress": [addr] };
    web3.trace.filter(filter, function (err, tx) {
      if (err || !tx) {
        console.error("TraceWeb3 error :" + err)
        res.write(JSON.stringify({ "error": true }));
      } else {
        res.write(JSON.stringify(filterTrace(tx)));
      }
      res.end();
    })
  } else if ("addr" in req.body) {
    var addr = req.body.addr.toLowerCase();
    var options = req.body.options;

    var addrData = {};

    if (options.indexOf('checksummedAddr') > -1) {
      addrData['checksummedAddr'] = web3.toChecksumAddress(addr);
    }
    if (options.indexOf("balance") > -1) {
      try {
        addrData["balance"] = web3.eth.getBalance(addr);
        addrData["balance"] = etherUnits.toEther(addrData["balance"], 'wei');
      } catch (err) {
        console.error("AddrWeb3 error :" + err);
        addrData = { "error": true };
      }
    }
    if (options.indexOf("count") > -1) {
      try {
        addrData["count"] = web3.eth.getTransactionCount(addr);
      } catch (err) {
        console.error("AddrWeb3 error :" + err);
        addrData = { "error": true };
      }
    }
    if (options.indexOf("bytecode") > -1) {
      try {
        addrData["bytecode"] = web3.eth.getCode(addr);
        if (addrData["bytecode"].length > 2)
          addrData["isContract"] = true;
        else
          addrData["isContract"] = false;
      } catch (err) {
        console.error("AddrWeb3 error :" + err);
        addrData = { "error": true };
      }
    }

    res.write(JSON.stringify(addrData));
    res.end();


  } else if ("block" in req.body) {
    var blockNum = parseInt(req.body.block);

    web3.eth.getBlock(blockNum, function (err, block) {
      if (err || !block) {
        console.error("BlockWeb3 error :" + err)
        res.write(JSON.stringify({ "error": true }));
      } else {
        res.write(JSON.stringify(filterBlocks(block)));
      }
      res.end();
    });

  } else {
    console.error("Invalid Request: " + action)
    res.status(400).send();
  }

};

exports.eth = web3.eth;
