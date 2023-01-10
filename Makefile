COMPOSE=cd docker/ && docker-compose

export DOCKER_UID=$(shell id -u)
export DOCKER_GID=$(shell id -g)

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

### Global

start: ## Start project
	make docker-up

stop: ## Stop project
	make docker-down

.PHONY: start stop

### Docker

docker-up: ## Start or restart all the docker services
	$(COMPOSE) up -d

docker-down:  ## Stop all the docker services
	$(COMPOSE) down

docker-rebuild: ## Rebuild docker environment
	$(COMPOSE) up -d --build --force-recreate

docker-start: ## Start existing docker environment
	$(COMPOSE) start

docker-stop: ## Stop existing docker environment
	$(COMPOSE) stop

docker-bash: ## Launch a shell into the container go
	$(COMPOSE) exec --user=$(DOCKER_UID):$(DOCKER_GID) golang bash

docker-log: ## Containers logs
	$(COMPOSE) logs

.PHONY: docker-log docker-up docker-down docker-rebuild docker-start docker-stop docker-bash

## Tweak to remove not rules message
%:
	@:
