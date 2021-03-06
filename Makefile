.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## build application dockers
	docker-compose build

run: ## run application dockers
	docker-compose up -d

stop: ## stop application dockers
	docker-compose stop

generate: ## generate dummy cats to DB using application functionality
	docker-compose run --rm --entrypoint "/app/main generate" crud ./...
