BINARY_NAME := server

PROTO_FILE := api/order.proto
PROTO_OUTPUT_DIR := pkg/api/test

GO := go
PROTOC := protoc
GOLANGCI_LINT := golangci-lint

.PHONY: all build run generate lint test clean help

all: build

build:
	$(GO) build -o $(BINARY_NAME) ./cmd

run: build
	./$(BINARY_NAME)

generate:
	@mkdir -p $(PROTO_OUTPUT_DIR)
	$(PROTOC) \
		--go_out=$(PROTO_OUTPUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUTPUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

lint:
	$(GOLANGCI_LINT) run

test:
	$(GO) test -v ./...

clean:
	rm -f $(BINARY_NAME)
	# При желании можно также удалить сгенерированные файлы:
	# rm -f $(PROTO_OUTPUT_DIR)/*.pb.go

help:
	@echo "Доступные команды:"
	@echo "  build     — собрать бинарный файл"
	@echo "  run       — собрать и запустить сервер"
	@echo "  generate  — пересоздать gRPC-код из order.proto"
	@echo "  lint      — проверить код с помощью golangci-lint"
	@echo "  test      — запустить все тесты"
	@echo "  clean     — удалить бинарник"
	@echo "  help      — показать эту справку"