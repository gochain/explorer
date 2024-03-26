.PHONY: docker test-backend admin server grabber backend frontend build generate

docker:
	docker build -t gochain/explorer .

test-backend:
	go test ./...

server:
	cd server && go build

grabber:
	cd grabber && go build

admin:
	cd admin && go build

backend: server grabber admin

frontend:	
	cd front && npm i
	# npm postintall not working in root user
	cd front && node patch.js
	rm -rf front/dist/explorer
	cd front && npm rebuild node-sass
	cd front && ./node_modules/@angular/cli/bin/ng build --prod
	cp -r front/dist .

build: backend frontend

generate:
	cd contracts && web3 contract build ERC721.sol && web3 contract build ERC20.sol && web3 contract build Upgradeable.sol
	cd server/tokens && abigen --lang go --abi ../../contracts/ERC20.abi --bin ../../contracts/ERC20.bin --pkg tokens --type ERC20 --out erc20.go
	cd server/tokens && abigen --lang go --abi ../../contracts/ERC721.abi --bin ../../contracts/ERC721.bin --pkg tokens --type ERC721 --out erc721.go
	cd server/tokens && abigen --lang go --abi ../../contracts/Upgradeable.abi --bin ../../contracts/Upgradeable.bin --pkg tokens --type Upgradeable --out upgradeable.go