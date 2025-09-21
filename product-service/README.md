# Product Services

## Install Dependency

```
cd product-service

go mod tidy

cp env.sample .env
```

## DB Migration

Create database via SQL

```
CREATE DATABASE product_service
```

Migrate table to database

```
make migrate MIGRATE_ARGS=up MYSQL_USER=root MYSQL_PASS=rootpw MYSQL_DB=product_service
```

Here's the sample command for create database migration file

```
make generate-db-migration MIGRATE_MODULE=product MIGRATE_NAME=create_table_product

```

## Running Service

```
go run cmd/gateway/main.go
```

## Database

### Table: products

```
id              bigint (primary key)
name            varchar(255)
price           decimal(15,3)
crated_at       timestamp
updated_at      timestamp
```

```
index :
- name
```

### Sample Insert Table

```
INSERT INTO `product_service`.`products` (`name`, `price`) VALUES ('Lorem Ipsum', 10000);
```

## INTERNAL API

### Product

```
URL: GET /check-products

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
    "products": [
		{
			"id": "1",
			"name": "Lorem Ipsum",
			"price": "10000"
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

## PUBLIC API

### Product

Called Internal Service:

- Shop
- Warehouse Stock

```
URL: GET /products

Parameters:
page_num = int # Default = 1
page_size = int # Default = 10
name = string
```

```json
Http Status: 200
Response:
{
    "products": [
		{
			"id": "1",
			"name": "Lorem Ipsum",
			"price": "10000",
            "total_stock": 10,
            "shops": [
                {
                    "id": "1",
                    "name": "Lorem Shop",
                    "total_stock": 5,
                    "warehouses": [
                        {
                            "warehouse_id": "1",
                            "warehouse_name": "Lorem Warehouse",
                            "warehouse_stock_id": "1",
                            "warehouse_stock": 3
                        }
                    ]
                }
            ]
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
