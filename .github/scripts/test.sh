#bin/bash

#run unit test
go test ./... -v

#run integration test
docker compose up -d

TEST_INTEGRATION=TRUE go test ./... -v 

docker compose down 
