SHELL:=/bin/bash
NAME =gennaker
SOURCES :=		$(shell find . -name "*.go")
DOCKER_IMAGE ?=		vgheri/$(NAME)
DOCKER_RUN_OPTS ?=	--rm
CURRENT_BRANCH_GIT ?= $(shell git symbolic-ref --short HEAD)
CURRENT_BRANCH ?= $(shell echo $(CURRENT_BRANCH_GIT) | tr / -)
TAG ?=		$(shell if [ $(CURRENT_BRANCH) == "master" ]; then echo "latest"; else echo $(CURRENT_BRANCH) | tr / -; fi)
GIT_SHA ?= $(shell git rev-parse HEAD)
BUILD_NUMBER ?= $(shell hostname)
DOCKER_COMPOSE_CMD ?= docker-compose -p $(TAG)
DOCKER_COMPOSE_EXEC_CMD ?= $(DOCKER_COMPOSE_CMD) exec -T

.PHONY: build
build: install

.PHONY: install
install:
	go install

.PHONY: run
run: install
	${GOPATH}/bin/$(NAME)

.PHONY: docker.build
docker.build:
	docker build -t $(DOCKER_IMAGE):$(GIT_SHA) .

.PHONY: docker.test
docker.test:
	$(DOCKER_COMPOSE_CMD) build --pull
	$(DOCKER_COMPOSE_CMD) up -d
	$(DOCKER_COMPOSE_EXEC_CMD) $(NAME) make test
	# $(DOCKER_COMPOSE_EXEC_CMD) $(NAME) integration
	$(DOCKER_COMPOSE_CMD) down

.PHONY: docker.push
docker.push: docker.build
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD)
	docker tag $(DOCKER_IMAGE):$(GIT_SHA) $(DOCKER_IMAGE):$(TAG)
	docker tag $(DOCKER_IMAGE):$(TAG) $(DOCKER_IMAGE):build-$(BUILD_NUMBER)
	docker push $(DOCKER_IMAGE):$(TAG)
	docker push $(DOCKER_IMAGE):build-$(BUILD_NUMBER)

.PHONY: docker.run
docker.run:
	$(DOCKER_COMPOSE_CMD) build --pull
	$(DOCKER_COMPOSE_CMD) up -d

.PHONY: docker.stop
docker.stop:
	$(DOCKER_COMPOSE_CMD) down

.PHONY: test
test:
	go vet $(shell go list ./... | grep -v /vendor/)
	go test -v -p 1 -race $(shell go list ./... | grep -v /vendor/)
