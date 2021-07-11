#bin/bash

#run unit test
TEST_INTEGRATION= go test ./... -v

#run integration test
docker compose up -d

TEST_INTEGRATION=TRUE go test ./... -v 

docker compose down 
