include .env

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o ./users-service ./cmd/users/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o ./migrator-service ./cmd/migrator/main.go

.PHONY: prepare
	go mod download

.PHONY: users-service
users--service:
	./users-service

.PHONY: migrator-service
migrator-service:
	./migrator-service

.PHONY: clean
clean:
	rm ./gofemart-service
	rm ./migrator-service

.PHONY: lint
lint:
	golangci-lint run ./... --fix

.PHONY: swag-generate
swag-generate:
	swag init -g  ./cmd/users/main.go --parseDependency --parseInternal

.PHONY: migrate_all_up migrate_all_down migrate_force migrate_version migrate_up migrate_down
migrate_all_up:
	migrate -database ${INFRA__POSTGRES__CONN_STR} -path /migrations up

migrate_all_down:
	migrate -database ${INFRA__POSTGRES__CONN_STR} -path /migrations down

migrate_force:
	migrate -database ${INFRA__POSTGRES__CONN_STR} -path /migrations force 1

migrate_version:
	migrate -database ${INFRA__POSTGRES__CONN_STR} -path /migrations version

migrate_up:
	migrate -database ${INFRA__POSTGRES__CONN_STR} -path /migrations up 1

migrate_down:
	migrate -database ${INFRA__POSTGRES__CONN_STR} -path /migrations down 1

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate_create name=your_migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./migrations -seq $(name)


.PHONY: start-app
start-app:
	docker-compose up --build -d

.PHONY: stop_app
stop-app:
	docker-compose down