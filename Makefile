.PHONY: test docker
docker:
	docker build -t gochain/explorer .

test:
	npm install
	./run_tests.sh
