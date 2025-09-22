# commerce-exercise

## Infrastructure

1. infrastructure : docker compose of mysql
2. service : docker compose for deployment services

## Development Services

1. user_service
2. shop_service
3. product_service
4. warehouse_service
5. order_service

## Running Docker Cron

```
docker run -it \
 --add-host=host.docker.internal:host-gateway \
 commerce-exercise-order-service/cron/expired-order:latest
```
