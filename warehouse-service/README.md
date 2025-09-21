# Warehouse Services

## Install Dependency

```
cd warehouse-service

go mod tidy

cp env.sample .env
```

## DB Migration

Create database via SQL

```
CREATE DATABASE warehouse_service
```

Migrate table to database

```
make migrate MIGRATE_ARGS=up MYSQL_USER=root MYSQL_PASS=rootpw MYSQL_DB=warehouse_service
```

Here's the sample command for create database migration file

```
make generate-db-migration MIGRATE_MODULE=warehouse MIGRATE_NAME=create_table_warehouse
make generate-db-migration MIGRATE_MODULE=warehouse MIGRATE_NAME=create_table_warehouse_stock

```

## Running Service

```
go run cmd/gateway/main.go
```

## Database

### Table: warehouses

```
id              bigint (primary key)
shop_id         bigint
name            varchar(255)
active          boolean
crated_at       timestamp
updated_at      timestamp
```

```
index :
- shop_id
```

### Table: warehouse_stocks

```
id              bigint (primary key)
warehouse_id    bigint
product_id      bigint
stock           int
crated_at       timestamp
updated_at      timestamp
```

```
unique index :
- warehouse_id, product_id

index:
- product_id, warehouse_id
```

### Sample Insert Table

```
INSERT INTO `warehouse_service`.`warehouses` (`shop_id`, `name`, `active`) VALUES (1, 'Lorem Ipsum Warehouse', true);
```

## INTERNAL API

### Stock

```
URL: GET /active-stocks

Authorization: Basic Auth

Parameters:
product_ids = array of int (required)
```

```json
Http Status: 200
Response:
{
    "warehouses": [
        {
            "id": "1",
            "shop_id": "1"
        }
    ],
    "warehouse_stocks": [
		{
			"id": "1",
            "warehouse_id": "1",
            "product_id": "1",
			"stock": 10
		}
	]
    "meta": {
        "http_status_code": 200
    }
}
```

### Adjustment Stock

```
URL: POST /adjustment-stocks

Authorization: Basic Auth
```

```json
Request:
{
    "warehouse_stocks": [
		{
            "warehouse_id": "1",
            "product_id": "1",
			"stock": 10,
		},
        {
            "warehouse_id": "1",
            "product_id": "2",
			"stock": -2,
		}
	]
}
```

```json
Http Status: 200
Response:
{
    "message": "Success adjustment stock",
    "meta": {
        "http_status_code": 200,
    }
}
```

### Transfer Stock

```
URL: POST /transfer-stocks

Authorization: Basic Auth
```

```json
Request:
{
    "original_warehouse_id": "1",
    "destination_warehouse_id": "2",
    "products": [
		{
            "product_id": "1",
			"stock": 10,
		},
        {
            "product_id": "2",
			"stock": 2,
		}
	]
}
```

```json
Http Status: 200
Response:
{
    "message": "Success transfer stock",
    "meta": {
        "http_status_code": 200,
    }
}
```

### Warehouse Activation

```
URL: POST /warehouse-actives

Authorization: Basic Auth
```

```json
Request:
{
    "warehouse_id": "1",
    "active": false
}
```

```json
Http Status: 200
Response:
{
    "message": "Success inactive/active warehouse",
    "meta": {
        "http_status_code": 200,
    }
}
```
