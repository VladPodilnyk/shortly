## help: prints help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## start: runs an application inside a docker container
.PHONY: start
start:
	docker build -f ./build/Dockerfile --tag url-shortener-app .
	docker run -p 4000:4000 url-shortener-app

## run-dev: run the cmd/api application
.PHONY: run-dev
run-dev:
	go run ./cmd/shortly

## audit: tidy dependencies and format
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...

## test: runs test
.PHONY: test
test:
	@echo 'Testing code...'
	go test -race -vet=off ./...