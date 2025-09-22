# Shop Services

## Install Dependency

```
cd shop-service

go mod tidy

cp env.sample .env
```

## DB Migration

Create database via SQL

```
CREATE DATABASE shop_service
```

Migrate table to database

```
make migrate MIGRATE_ARGS=up MYSQL_USER=root MYSQL_PASS=rootpw MYSQL_DB=shop_service
```

Here's the sample command for create database migration file

```
make generate-db-migration MIGRATE_MODULE=shop MIGRATE_NAME=create_table_shop

```

## Running Service

```
go run cmd/gateway/main.go
```

## Build Image

```
make compile
make build
```

## Database

### Table: shops

```
id              bigint (primary key)
name            varchar(255)
crated_at       timestamp
updated_at      timestamp
```

### Sample Insert Table

```
INSERT INTO `shop_service`.`shops` (`name`) VALUES ('Lorem Ipsum Store');
```

## INTERNAL API

### Shop

```
URL: GET /shops

Authorization: Basic Auth

Parameters:
page_num = int # Default = 1
page_size = int # Default = 10
ids = array of int
```

```json
Http Status: 200
Response:
{
    "shops": [
		{
			"id": "1",
			"name": "Lorem Ipsum",
		}
	]
    "meta": {
        "http_status_code": 200,
        "page_num": 1,
        "page_size": 10,
        "page_total": 5
    }
}
```
