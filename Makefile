up: 
	docker-compose up -d

down:
	docker-compose down

build:
	go build .

run:
	./todo-app

go:
	make build && make run

all:
	make up && make go

clean:
	rm todo-app && make down
