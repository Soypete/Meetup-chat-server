# SETUP special settings db locally
Some blogs i found useful: (go rest api with chi)[https://blog.logrocket.com/how-to-build-a-restful-api-with-docker-postgresql-and-go-chi/] and (postgres in docker)[https://betterprogramming.pub/how-to-run-postgresql-in-a-docker-container-on-macos-with-persistent-local-data-440e7172821e]


# Quickstart
## docker and postgres

1. Make a data directory for the pg version.
2. create docker network (hash id will be returned)
3. postgres container (hash id will be returned)

```sh
# 1
mkdir -p ~/Pgdata/postgres_data/12.0

# 2 
docker network create dev-network
## => 92ba569174278df6c5415c0c94fb0c59f4b1b8aafb103fe446bd9e45c5734686

# 3 :5432 may conflict with an already running postgres server...
#docker run --restart always --name IMAGE_NAME --net dev-network -v /Users/[YOUR_USERNAME]/[YOUR_DATA_DIRECTORY]:/var/lib/postgresql/data -p 5432:LOCAL_PORT -d -e POSTGRES_PASSWORD=[YOUR_PASSWORD] -e POSTGRES_DB=[YOUR_DATABASE] postgres:PG_VERSION
docker run --restart always --name pgsettings12 --net dev-network -v /Users/josh.bowlesgetweave.com/Pgdata/postgres_data/12.0:/var/lib/postgresql/data -p 5432:5432 -d -e POSTGRES_PASSWORD=settingsPass -e POSTGRES_DB=sync_settings postgres:12
## => 5c75a472349d9e7babb0bfe57450a9444c67952e1efaeca825853e8ec5ca152b
```
## check docker
Look in desktop dockerapp dashboard for `pgsettings12` or do `docker ps` in terminal.

```sh
docker ps
#CONTAINER ID   IMAGE         COMMAND                  CREATED         STATUS         PORTS                                       NAMES
#5c75a472349d   postgres:12   "docker-entrypoint.s…"   8 seconds ago   Up 7 seconds   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp   pgsettings12
```

## test connecting to the db
If you are running a local postgres server on 5432 then you should shut it down.

```sh
## -W will prompt for the POSTGRES_PASSWORD
psql -h localhost -U postgres -W
# list databases... you should see none
\l
```

## create the db

```sh
## will prompt for the POSTGRES_PASSWORD; no answer means success!
createdb -h localhost -p 5432 -U postgres sync_settings
```

## get migrate project and create migration files

```sh
## https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
brew install golang-migrate
```

```sh
## write the SQL for migrations as well as rollback; this creates up/down files
migrate create -ext sql -dir db/migrations -seq create_all_tables
```

## run migrations

```sh
export POSTGRESQL_URL="postgres://postgres:settingsPass@localhost:5432/sync_settings?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path db/migrations up
## => migrate -database ${POSTGRESQL_URL} -path db/migrations up
## => 1/u create_all_tables (189.489557ms)

psql -h localhost -U postgres -d sync_settings -W
\d
```

## Setup a postgres UI
Like any othe db connection just parse the url into the feilds needed to connect: host, port, user, password. Same with code and/or scripts.


# Other
## more about the docker command

* `--restart always` will restart this container any time Docker is started, such as for a laptop reboot or if Docker gets shut down and started again. Leave this parameter out if you want to start your own containers every time.
* `--name pgsettings12` assigns the name "pgSettings12.0" to your container instance. This is the way it will appear in the Docker dashboard. Adding the version number to the name is handy if you’re planning on running multiple versions in containers on your Mac.
* `--network dev-network` will join this container to your local Docker network that you created.
* `-v Users/josh.bowlesgetweave.com/Pgdata/postgres_data/12.0:/var/lib/postgresql/data` will bind/mount data folder inside the container volume (/var/lib/postgresql) to the local directory on laptop; data will persist so we can restart.
* `-p 5432:5432` will bind the Postgres port of the container (5432) to the same port on your Mac. This will make things easier in treating this container as a “localhost” in Postgres.
* `-d` will run this container in detached mode so that it runs in the background.
* `-e POSTGRES_PASSWORD=[YOUR_PASSWORD]` and `-e POSTGRES_DB=[YOUR_DATABASE]` sets an environment variable (in this case, the PostgreSQL root password inside the container).
* `postgres:12` indicates that the official DockerHub Postgres version tag 13.2 is the one to install. Simply replace the :13.2 with a different version if you’re creating a new instance (e.g. :12.0) and make sure you’ve created a data folder on your Mac that matches so you can bind to it.
