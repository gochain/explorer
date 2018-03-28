.PHONY: test docker
docker:
	docker build -t gochain/explorer .

test:
	./run_tests.sh
