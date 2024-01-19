
build:
	go build -o ./bin/server ./cmd/server

run: build
	./bin/server

docker_clean:
	docker-compose down -v

docker_start:
	docker-compose up -d db
	docker-compose up --build setup
	docker-compose up --build app