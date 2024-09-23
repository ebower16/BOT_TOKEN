APP_NAME=botus
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_COMPOSE_FILE=docker-compose.yml

.PHONY: all build run test docker up down clean

all: build

build:
	@echo "Building the application..."
	go build -o $(APP_NAME) ./cmd

run: build
	@echo "Running the application..."
	./$(APP_NAME)

test:
	@echo "Running tests..."
	go test ./...

docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

up:
	@echo "Starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

down:
	@echo "Stopping Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	docker rmi $(DOCKER_IMAGE) || true
