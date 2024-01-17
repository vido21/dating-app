.PHONY: init build run test coverage mock

init:
	cp .env.example .env
	docker-compose up -d
	go mod vendor

build:
	go build -v

run:
	go run main.go

test:
	go test -coverprofile cover.out ./...

cover:
	go tool cover -html=cover.out -o cover.html

mock:
	mockery --all

generate-coverage:
	chmod +x generate-test-coverage.sh
	./generate-test-coverage.sh
