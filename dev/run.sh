#bin/bash

go mod tidy
go mod download

#TODO: check if container is running

docker run --restart always --name pgtwitch --net d
ev-network -v /Users/miriahpeterson/Pgdata/twitch_db:/var/lib/postgresql/data -p 5432:5432 -d
 -e POSTGRES_PASSWORD=postgres postgres:12

go run main.go
