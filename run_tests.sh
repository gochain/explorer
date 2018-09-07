set -e
docker run --name test_mongo -d -p 8545:8545 -p 8080:8080 -p 27017:27017 mongo
varA=`docker ps --no-trunc -q | cut -c 1-12`
docker build -t gochain/explorer:test_ci .
docker run --name test_explorer_grabber -d --network="container:$varA" gochain/explorer:test_ci grabber -u https://testnet-rpc.gochain.io -s 10
docker run --name test_explorer_server -d --network="container:$varA" gochain/explorer:test_ci server -d /explorer/ -u https://testnet-rpc.gochain.io
sleep 20 # let's wait until server start
# docker exec test_explorer npm test
echo "Docker logs for grabber"
docker logs test_explorer_grabber
echo "Docker logs for server"
docker logs test_explorer_server
echo "Trying curl"
docker run --rm --network="container:$varA" byrnedo/alpine-curl -f http://localhost:8080/
docker run --rm --network="container:$varA" byrnedo/alpine-curl -f http://localhost:8080/api/blocks/10 
docker run --rm --network="container:$varA" byrnedo/alpine-curl -f http://localhost:8080/api/blocks/10/transactions
for imagename in test_web3 test_mongo test_explorer_grabber test_explorer_server; do
    echo "* Trying to stop ${imagename}..."
    [ -z $(docker ps -a -f "name=$imagename" --format '{{.Names}}') ] || docker rm -f $imagename
done
