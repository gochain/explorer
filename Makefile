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
	cd front && npm i
	# npm postintall not working in root user
	cd front && node patch.js
	rm -rf front/dist/explorer
	cd front && npm rebuild node-sass
	cd front && ./node_modules/@angular/cli/bin/ng build --prod
	cp -r front/dist .

build: backend frontend
