DOCKER_IMAGE ?=	vpcorp/gennaker-pgsql-ssl

.PHONY: build
build: docker.build

.PHONY: docker.build
docker.build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker.run
docker.run:
	docker run -p 5432:5432 -d -e POSTGRES_DB=$(DBNAME) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) vpcorp/gennaker-pgsql-ssl
