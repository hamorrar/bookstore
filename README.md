# bookstore

## To Set Up Migration (up and down)

### Install Go Migrate
```bash
curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
sudo apt-get update
sudo apt-get install -y migrate
```
### Set Database Environment Variable
```bash
export DB_URL="postgres://postgres:postgres@localhost:5432/bookstore?sslmode=disable"
```
### Create migration files
```bash
migrate create -ext sql -dir migrations -seq create_users_table
migrate create -ext sql -dir migrations -seq create_books_table
migrate create -ext sql -dir migrations -seq create_orders_table
```

### Apply Migration Up
```bash
migrate -path migrations -database "postgresql://username:password@localhost:5432/database_name?sslmode=disable" up
```

## Port Check up
```bash
sudo ss -tulnp | grep :<PORT_NUMBER>
```