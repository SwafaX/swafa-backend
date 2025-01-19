up: 
	docker-compose up -d

down:
	docker-compose down

build:
	go build .

run:
	./swafa-backend

go:
	make build && make run

all:
	make up && make go

clean:
	make down && rm swafa-backend
