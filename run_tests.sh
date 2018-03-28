docker run -d -p 8545:8545 -p 3000:3000 -p 27017:27017 trufflesuite/ganache-cli
varA=`docker ps --no-trunc -q | cut -c 1-12`
docker run -d --network="container:$varA" mongo
docker build -t gochain/explorer:test_ci .
docker run -d --network="container:$varA" -e RPC_HOST=localhost  gochain/explorer:test_ci
sleep 5 # let's wait until server start
npm test
docker rm -f $(docker ps -a -q)
