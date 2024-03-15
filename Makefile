.SILENT:
.DEFAULT_GOAL := fast-run

.PHONY: build
build:
	go build -o ./build/calc ./cmd/calc

.PHONY: run
run:
	./build/calc

.PHONY: fast-run
fast-run:
	go run ./cmd/calc/main.go
