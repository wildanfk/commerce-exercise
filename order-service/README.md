# Order Services

## Install Dependency

```
cd order-service

go mod tidy

cp env.sample .env
```

## DB Migration

Create database via SQL

```
CREATE DATABASE order_service
```

Migrate table to database

```
make migrate MIGRATE_ARGS=up MYSQL_USER=root MYSQL_PASS=rootpw MYSQL_DB=order_service
```

Here's the sample command for create database migration file

```
make generate-db-migration MIGRATE_MODULE=order MIGRATE_NAME=create_table_order
make generate-db-migration MIGRATE_MODULE=order MIGRATE_NAME=create_table_order_detail

```

## Running Service

```
go run cmd/gateway/main.go
```

## Running Cron Service

```
go run cmd/cron/expired-order/main.go

Called Internal Service:

- Warehouse Stock
```

## Build Image

```
make compile
make build
```

## Database

### Table: orders

```
id              bigint (primary key)
user_id         bigint
shop_id         bigint
state           tinyint
total_stock     int
total_price     decimal(15,3)
expired_at      datetime
crated_at       timestamp
updated_at      timestamp
```

```
index :
- user_id
- shop_id
- state, expired_at
```

### Table: order_details

```
id              bigint (primary key)
order_id        bigint
product_id      bigint
warehouse_id    bigint
stock           int
price           decimal(15,3)
crated_at       timestamp
updated_at      timestamp
```

```
index:
- order_id
```

## PUBLIC API

### Order Checkout

Called Internal Service:

- Product
- Warehouse Stock

```
URL: POST /checkout-orders

Authorization: User Auth
```

```json
Request:
{
    "shop_id": "1",
    "products": [
		{
            "id": "1",
            "warehouse_id": "1",
			"stock": 2
		}
	]
}
```

```json
Http Status: 201
Response:
{
    "message": "Success create order",
    "meta": {
        "http_status_code": 201,
    }
}
```
