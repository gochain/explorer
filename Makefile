.PHONY: dep server grabber build frontend backend docker release install test deploy

docker:
	docker build -t gochain/explorer:go-test .

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
	cd server &&  go build -v
	cd grabber && go build -v

frontend:
	npm i
	# npm postintall not working in root user
	node patch.js
	rm -rf dist/explorer
	npm rebuild node-sass
	ng build --prod
