PROJECT_NAME = PR-Service

run: down build up

stop:
	docker stop $$(docker ps -a -q)

remove:
	docker rm $$(docker ps -a -q)

gen:
	oapi-codegen --config=api/v1/oapi-codegen.yaml api/v1/openapi.yml

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

clean:
	docker-compose down -v
	docker system prune -f

migrate:
	docker-compose run --rm migrate

dev:
	go run cmd/pr-service/main.go

lint:
	golangci-lint run ./... -v

lint-fix:
	golangci-lint run ./... --fix -v

install-lint:
	@which golangci-lint > /dev/null || (echo "install golangci-lint" && \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin)
