.PHONY: run build test migrate migrate-down clean

run:
	docker-compose up --build

build:
	go build -o main ./cmd/httpserver

test:
	go test ./... -cover -v

migrate:
	docker-compose exec app migrate -path /migrate -database "postgres://$$DB_USER:$$DB_PASSWORD@db:5432/$$DB_NAME?sslmode=disable" up

migrate-down:
	docker-compose exec app migrate -path /migrate -database "postgres://$$DB_USER:$$DB_PASSWORD@db:5432/$$DB_NAME?sslmode=disable" down

clean:
	docker-compose down -v
	docker system prune -f
	docker volume prune -f
