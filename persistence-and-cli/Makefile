BINARY := blockchain

all: build test

build:
	@echo "==> Go build"
	@go build -o $(BINARY)

depends:
	@go get -u github.com/boltdb/bolt

test:
	@echo "==> 开始运行"
	@./$(BINARY) print
	@./$(BINARY) add -data "给彭爽1比特币"
	@./$(BINARY) add -data "彭爽花了0.01比特币买咖啡"
	@./$(BINARY) print

.PHONY: build depends test
