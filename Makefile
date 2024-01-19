GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)
FRONTEND:= "frontend"
MOCKS_DESTINATION=mocks
.PHONY: mocks

## Docker:
docker-app-up: ## Create and run app containers
	docker compose --file docker-compose/app/docker-compose.yml up -d --force-recreate --build

docker-app-down: ## Stop and remove app containers
	docker compose --file docker-compose/app/docker-compose.yml down -v

docker-dev-db-up: ## Create and run dev container with db
	docker compose --file docker-compose/dev/docker-compose.yml up -d --force-recreate

docker-dev-db-down: ## Stop and remove dev container with db
	docker compose --file docker-compose/dev/docker-compose.yml down -v

docker-test-db-up: ## Create and run test container with db
	docker compose --file docker-compose/test/docker-compose.yml up -d --force-recreate

docker-test-db-down: ## Stop and remove test container with db
	docker compose --file docker-compose/test/docker-compose.yml down -v

fmt:
	gofmt -s -w .
	goimports -w .

mocks: $(shell grep -lrP --include='*.go' --exclude='*test.go' 'type\s+\w+\s+interface\s*' ./)
	@echo "Generating mocks"
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/`echo $${file#*internal/}`; done

swag:
	swag fmt
	swag init -g ./internal/app/rest/handlers.go

build: swag
	go build -v -o qon ./cmd/main.go

## Test:
test: ## Run tests
	@docker compose --file docker-compose/test/docker-compose.yml up -d --force-recreate
	@go test -count=1 -v ./...
	@docker compose --file docker-compose/test/docker-compose.yml down -v

test-coverage: ## run test and show coverage
	@docker compose --file docker-compose/test/docker-compose.yml up -d
	@echo "Package test coverage:"
	@go test -count=1 -coverpkg=./internal/... -coverprofile=coverage.out ./...
	@echo "\n\n"
	@echo "Separate files test coverage:"
	@go tool cover -func coverage.out
	@docker compose --file docker-compose/test/docker-compose.yml down
	@timeout 5 echo
	@rm coverage.out

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
