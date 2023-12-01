GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Docker:

docker-dev-db-up: ## Create and run dev container with db
	docker compose --file docker-compose/dev/docker-compose.yml up -d --force-recreate

docker-dev-db-down: ## Stop and remove dev container with db
	docker compose --file docker-compose/dev/docker-compose.yml down

## Test:
test: ## Run tests
	@docker compose --file docker-compose/test/docker-compose.yml up -d
	@go test -count=1 -v ./...
	@docker compose --file docker-compose/test/docker-compose.yml down

## Info:
info: ## Show help information
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

.DEFAULT_GOAL = info
