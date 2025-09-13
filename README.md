# bookstore

## Set Up

### Download source code with SSH
```bash
git clone git@github.com:hamorrar/bookstore.git; cd bookstore
```

### Install Go Dependencies
```bash
go mod tidy
```

### Install Go Migrate
```bash
curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
sudo apt-get update
sudo apt-get install -y migrate
```

### How to set up .env files
- Put ``.env`` file in ``/``.
- Define DB_NAME, SECRET_KEY, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, PORT, DB_URL, DEFAULT_DB_URL, DB_DSN, DEFAULT_DB_DSN.

### Documentation
```bash
swag init --dir cmd/api --parseDependency --parseInternal --parseDepth 1
```
Might need to ``export PATH="$(go env GOPATH)/bin:$PATH"`` first.

### Log into DB in a terminal
```bash
psql -U postgres -h localhost
```
Add ``-d <DATABASE_NAME>`` to connect to a specific database.

### Create migration files
```bash
migrate create -ext sql -dir migrations -seq create_users_table
migrate create -ext sql -dir migrations -seq create_books_table
migrate create -ext sql -dir migrations -seq create_orders_table
```

### Apply Migration Up
can replace up/down as needed
```bash
migrate -path migrations -database $DB_URL up
```
or
```bash
go run ./cmd/migrate/main.go up
```

### Start Go Server
```bash
go run ./cmd/api
```

## To Run from the root directory
1. Apply up migrations as above
1. start the go server
1. send client requests with ``curl``

## Testing
### Go Tests
```bash
go test ./...
```
Flags:
- ``-v`` flag for verbose mode
- ``-failfast`` to stop on the first failure
- ``-cover`` to show test coverage percentage.

## Misc

### Port Check up
```bash
sudo ss -tulnp | grep :<PORT_NUMBER>
```

### Migration Errors
#### Version or Dirty table
```bash
migrate -path "cmd/migrate/migrations/" -database $DB_URL force 1
```

#### Drop migration table in psql terminal
```sql
drop table schema_migrations;
```

### Check migration table version/dirty
```sql
select * from schema_migrations;
```

## Ideas
- create a simulation of a big processing jobs by using a timer to wait for a random 1-5 second wait per "task" for 10 tasks to practice making a "/api/check/:jobid" endpoint to get the status of my job.