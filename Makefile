## help: prints help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: builds an application and stores a binary in ./bin
.PHONY: build
build:
	bash scripts/build_styles.sh
	go build -o ./bin/shortly ./cmd/shortly

## build-styles: builds styles file with tailwindcss
.PHONY: build-styles
build-styles:
	bash scripts/build_styles.sh

## build-container: builds a docker container for the app
.PHONY: build-container
build-container:
	docker build -f ./deploy/Dockerfile --tag shortly .

## run-dev: runs the cmd/api application
.PHONY: run-dev
run-dev:
	bash scripts/run_dev.sh

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

## tests: runs tests
.PHONY: tests
tests:
	@echo 'Testing code...'
	docker compose up -d
	go test -race -vet=off ./...

## cleanup: cleans up dev artifacts
.PHONY: cleanup
cleanup:
	rm ./cmd/shortly/public/output.css
	rm -rf ./bin
	rm -rf ./build/