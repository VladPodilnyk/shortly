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

.PHONY: compile
compile:
	go build -o ./bin/shortly ./cmd/shortly

.PHONY: build-styles
build-styles:
	bash scripts/build_styles.sh

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

## test: runs test
.PHONY: test
test:
	@echo 'Testing code...'
	docker-compose up -d
	go test -race -vet=off ./...

.PHONY: cleanup
cleanup:
	rm ./cmd/shortly/public/output.css
	rm -rf ./bin
	rm -rf ./build/