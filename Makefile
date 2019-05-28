.PHONY: docker test-backend server grabber backend frontend build

docker:
	docker build -t gochain/explorer .

test-backend:
	go test ./...

server:
	cd server && go build -v

grabber:
	cd grabber && go build -v

backend: server grabber

frontend:
	npm i
	# npm postintall not working in root user
	node patch.js
	rm -rf dist/explorer
	npm rebuild node-sass
	./node_modules/@angular/cli/bin/ng build --prod

build: backend frontend
