SHELL=/bin/bash -o pipefail

.PHONY: build
build:
	go build -o ./bin/current-rate-server ./cmd/.

.PHONY: run
run: build
	./bin/current-rate-server $(ARGS)

.PHONY: build-and-run
build-and-run: build run