.PHONY: dep server grabber build frontend backend docker release install test deploy

docker:
	docker build -t gochain/explorer .

# test:
# 	npm install
# 	./run_tests.sh

# release: docker
# 	./release.sh

server: build
	cd server && ./server

grabber: buildback
	cd grabber && ./grabber

build: backend frontend

backend:
	cd server && dep ensure --vendor-only && go build	
	cd grabber &&  dep ensure --vendor-only && go build

frontend:
	npm i
	rm -rf dist/explorer
	ng build --aot	