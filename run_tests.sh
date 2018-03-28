docker run --name test_web3 -d -p 8545:8545 -p 3000:3000 -p 27017:27017 trufflesuite/ganache-cli
varA=`docker ps --no-trunc -q | cut -c 1-12`
docker run --name test_mongo -d --network="container:$varA" mongo
docker build -t gochain/explorer:test_ci .
docker run --name test_explorer -d --network="container:$varA" -e RPC_HOST=localhost  gochain/explorer:test_ci
sleep 5 # let's wait until server start
npm test
for imagename in test_web3 test_mongo test_explorer; do
    echo "* Trying to stop ${imagename}..."
    [ -z $(docker ps -a -f "name=$imagename" --format '{{.Names}}') ] || docker rm -f $imagename
done
