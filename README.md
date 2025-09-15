# bookstore

## Summary
A RESTful API to manage a bookstore with basic CRUD functionality for users, books, and orders. Role-based access for admin and customers. The backend service is written in Go, PostgreSQL for data storage, and Gin framework for routing and middleware. User authentication uses JWTs stored in secure cookies and middleware validates tokens and enforces role-based control. It includes a simple CI/CD pipeline with GitHub Actions to build and test the service. The API supports pagination Swagger is used for API documentation and can be seen at [localhost:8080/swagger](http://localhost:8080/swagger) when the server is running and in [./docs](./docs) for the json/yaml files.

## Ideas/Plans
- Create a simulation of a big processing jobs to practice making an endpoint to check on job status.
- Use Docker to containerize the API and database.
- More robust testing.
- Look into how Terraform can be applied to the project for practice with infrastructure as code. Deploy product to AWS with Terraform.
- Improve CI/CD pipeline after building and testing jobs finish. Add linting.

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

### How to set up environment variables
- Put ``.env`` file in ``/``.
- Define DB_NAME, SECRET_KEY, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, PORT, DB_URL, DB_DSN.
- Put the SECRET_KEY and DB_PASSWORD in GitHub Secrets to reference in CI/CD workflow.

### Swagger set up
```bash
swag init --dir cmd/api --parseDependency --parseInternal --parseDepth 1
```
Might need to ``export PATH="$(go env GOPATH)/bin:$PATH"`` first.

### Log into DB in a terminal
```bash
psql -U postgres -h localhost
```
Add ``-d <DATABASE_NAME>`` to connect to a specific database. Ensure database is created before applying migrations below.

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
1. start the go server as above
1. send client requests with ``curl``. Examples below

### Example curl requests

To create an example admin user locally:
```bash
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "email": "user2@gmail.com",
  "password": "password2",
  "role": "Admin"
}' \
-w "\nHTTP Status: %{http_code}\n" \
http://localhost:8080/api/v1/auth/register
```
To login the example admin user locally:
```bash
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "email": "user2@gmail.com",
  "password": "password2"
}' \
-c cookies.txt \
-w "\nHTTP Status: %{http_code}\n" \
http://localhost:8080/api/v1/auth/login
```

To create a book with the admin user locally:
```bash
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "Author": "Author1",
  "Title": "Title1",
  "Price": 11
}' \
-b cookies.txt \
-w "\nHTTP Status: %{http_code}\n" \
http://localhost:8080/api/v1/books
```
``-b`` and ``-c`` flags are only necessary with ``curl`` to have a place to store the cookie locally. ``-b`` reads the cookie and ``-c`` reads the cookie from the specified file. The server creates and stores a cookie for each client.

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