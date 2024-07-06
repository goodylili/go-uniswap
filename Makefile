.PHONY: all build test clean docker

BINARY_NAME := tugbot
DOCKER_IMAGE_NAME := tugbot

all: build

build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) ./cmd

run:
	@echo "Running the application..."
	@go run ./cmd

env:
	@./load_env

test:
	go test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

push:
	@echo "Pushing to GitHub..."
	git add .
	@read -p "Enter commit message: " commit_msg; \
	git commit -m "$$commit_msg"; \
	git push