BUILD := "./build/anti_bruteforce"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

fmt:
	golangci-lint run --fix

build: swag build-app

build-app:
	go build -v -o $(BUILD) -ldflags "$(LDFLAGS)" .

run: build
	$(BUILD) serve --config "./configs/config.yaml"

build-img-app:
	docker build --no-cache \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_APP) \
		-f Dockerfile .

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
	$(BIN) migrate --config "./config.yaml"

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
	docker compose -f ./deployments/docker-compose.test.yaml up tests

down-test:
	docker compose -f ./deployments/docker-compose.test.yaml down --volumes

integration-tests:
	make build-img-tests
	make up-test
	make down-test

clear:
	make down-test
	make down
	make docker-clear

kub-info:
	kubectl get pods -A

helm-dep:
	helm dependency build anti-bruteforce-chart

helm-up-http:
	helm install anti-bruteforce-app chart --insecure-skip-tls-verify --namespace anti-bruteforce-app --create-namespace

helm-down-http:
	helm uninstall anti-bruteforce-app --namespace anti-bruteforce-app

.PHONY: build run build-img run-img version test lint down-test
