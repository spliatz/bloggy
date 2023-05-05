build:
	docker-compose build

start:
	docker-compose up

deps:
	go mod download & go install github.com/swaggo/swag/cmd/swag@latest

docs:
	rm -rf docs & swag init -g cmd/app/main.go
