.PHONY: dep server grabber build buildback buildfront docker release install test deploy

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

build: buildback buildfront

buildback:
	cd server && go get && go build	
	cd grabber && go get && go build
buildfront:
	npm i
	rm -rf dist/explorer
	ng build --aot	