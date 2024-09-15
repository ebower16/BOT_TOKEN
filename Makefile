# Имя исполняемого файла
BINARY_NAME=botus

GO_FILES=$(shell find ./cmd -name '*.go') $(shell find ./internal -name '*.go') $(shell find ./pkg -name '*.go')

# Цели Makefile
.PHONY: all
all: build

.PHONY: build
build:
	@echo "Сборка приложения..."
	go build -o $(BINARY_NAME) $(GO_FILES)

.PHONY: run
run: build
	@echo "Запуск приложения..."
	./$(BINARY_NAME)

.PHONY: test
test:
	@echo "Запуск тестов..."
	go test ./...

.PHONY: clean
clean:
	@echo "Очистка..."
	rm -f $(BINARY_NAME)

.PHONY: all-run
all-run: build run test