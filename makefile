include .env

BINARY_NAME=server
CMD_DIR=./cmd
PROTO_DIR=./proto
OUT_DIR=./pkg/api
PROTO_FILE=$(PROTO_DIR)/test.proto

GO=go
PROTOC=protoc

.PHONY: help
help:
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: proto-gen
proto-gen: check-deps
	@echo "Генерация gRPC-кода из $(PROTO_FILE)..."
	$(PROTOC) --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)
	@echo "Генерация завершена."

.PHONY: build
build:
	@echo "Сборка бинарника..."
	$(GO) build -o $(BINARY_NAME) $(CMD_DIR)
	@echo "Бинарник $(BINARY_NAME) собран."

.PHONY: run
run: build
	@echo "Запуск сервера..."
	./$(BINARY_NAME)

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	@echo "Бинарник удален."