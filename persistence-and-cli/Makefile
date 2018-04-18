BINARY := blockchain

all: build

build:
	@echo "==> Go build"
	@go build -o $(BINARY)

depends:
	@go get -u github.com/boltdb/bolt

.PHONY: build depends test
