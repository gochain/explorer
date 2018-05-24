.PHONY: test docker release
docker:
	docker build -t gochain/explorer .

test:
	npm install
	./run_tests.sh

release: docker
	./release.sh
