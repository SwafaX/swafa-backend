build:
	docker compose build

up: 
	docker-compose up -d

down:
	docker-compose down

clean:
	make down && rm swafa-backend

migrate:
	go run migrations/migrate.go
