.PHONY: dep run build docker release install test deploy

docker:
	docker build -t gochain/explorer .

test:
	npm install
	./run_tests.sh

release: docker
	./release.sh

run:
	ng serve --host 0.0.0.0

build:
	ng build --prod --aot

runprod: build
	ruby -run -e httpd ./dist/ -p 8080

deploy: build
	firebase deploy
