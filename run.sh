#bin/bash

go mod tidy
go mod download

docker-compose up --remove-orphans -d

go run main.go
