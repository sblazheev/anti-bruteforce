BUILD := build/anti_bruteforce
DOCKER_IMG_APP="anti_bruteforce:develop"
DOCKER_IMG_TESTS="tests:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.DEFAULT_GOAL: build

build: swag build-app

fmt:
	golangci-lint run --fix

build-app:
	go build -v -o $(BUILD) -ldflags "$(LDFLAGS)" .

run-cli: build
	$(BUILD) serve --config "./configs/config.yaml"

run: up

build-img-app:
	docker build --no-cache \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_APP) \
		-f Dockerfile .

build-img: build-img-app

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -count 100 ./... -tags !integrations

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.7.2

lint: install-lint-deps
	golangci-lint run ./...

migrate: build
	$(BUILD) migrate --config "./configs/config.yaml"

install-swag-deps:
	(which swag > /dev/null) || go install github.com/swaggo/swag/cmd/swag@latest

swag: install-swag-deps
	swag i -d internal/server/http/,internal/app/,internal/common/ -o internal/server/http/docs/ -g server.go

up:
	docker compose -f ./deployments/docker-compose.yaml up --force-recreate -d

down:
	docker compose -f ./deployments/docker-compose.yaml down

docker-clear:
	docker system prune --all --volumes -f

up-test:
	docker compose -f ./deployments/docker-compose.test.yaml up tests --force-recreate

down-test:
	docker compose -f ./deployments/docker-compose.test.yaml down --volumes

integration-tests:
	make up-test
	make down-test

clear:
	make down-test
	make down
	make docker-clear

.PHONY: build run build-img run-img version test lint down-test
