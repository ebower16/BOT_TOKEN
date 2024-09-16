
BINARY_NAME=botus


GO_FILES=$(shell find ./cmd -name '*.go') $(shell find ./internal -name '*.go')


.PHONY: all
all: build

.PHONY: build
build:
	@echo 
	go build -o $(BINARY_NAME) $(GO_FILES)

.PHONY: run
run: build
	@echo 
	./$(BINARY_NAME)

.PHONY: test
test:
	@echo 
	go test ./...

.PHONY: clean
clean:
	@echo 
	rm -f $(BINARY_NAME)

.PHONY: all-run
all-run: build run test
