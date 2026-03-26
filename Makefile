.PHONY: build run test bench fmt clean

BIN_DIR := bin
BIN := $(BIN_DIR)/firework

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN) ./cmd/firework/main.go

demo:
	go run ./cmd/firework/main.go -f ./demos/demo.csv $(ARGS)

run:
	go run ./cmd/firework/main.go $(ARGS)

test:
	go test ./...

bench:
	go test ./... -run '^$$' -bench . -benchmem

fmt:
	gofmt -w .

clean:
	rm -rf $(BIN_DIR)
