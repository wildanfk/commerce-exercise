# User Services

## Install Dependency

```
cd user-service

go mod tidy

cp env.sample .env
```

## DB Migration

Create database via SQL

```
CREATE DATABASE user_service
```

Migrate table to database

```
make migrate MIGRATE_ARGS=up MYSQL_USER=root MYSQL_PASS=rootpw MYSQL_DB=user_service
```

Here's the sample command for create database migration file

```
make generate-db-migration MIGRATE_MODULE=auth MIGRATE_NAME=create_table_user

```

## Running Service

```
go run cmd/gateway/main.go
```

## Database

### Table: users

```
id              bigint (primary key)
name            varchar(255)
email           varchar(255)
phone           varchar(100)
password        varchar(100)
crated_at       timestamp
updated_at      timestamp
```

```
index :
- email (unique)
- phone (unique)
```

### Sample Insert Table

```
INSERT INTO `user_service`.`users` (`name`, `email`, `phone`, `password`) VALUES ('Jhon Doe', 'jhon.doe@test.com', '+6281234567890', SHA2('test1234', 256));
```

## PUBLIC API

### Authentication

```
URL: POST /authentication
Content-Type: application/json
```

```json
Request:
{
	"username": "+6281234567890",
	"password": "jhondoe"
}
or
{
	"username": "jhon.doe@test.com",
	"password": "jhondoe"
}
```

```json
Http Status: 200
Response:
{
    "auth": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
    },
    "meta": {
        "http_status_code": 200
    }
}
```
