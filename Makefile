
APP_NAME=botus
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_COMPOSE_FILE=docker-compose.yml

.PHONY: all build run test docker up down clean

all: build


build:
	@echo
	go build -o $(APP_NAME) ./cmd


run: build
	@echo
	./$(APP_NAME)


test:
	@echo
	go test ./...


docker:
	@echo
	docker build -t $(DOCKER_IMAGE) .


up:
	@echo
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build


down:
	@echo
	docker-compose -f $(DOCKER_COMPOSE_FILE) down


clean:
	@echo
	rm -f $(APP_NAME)
	docker rmi $(DOCKER_IMAGE) || true
