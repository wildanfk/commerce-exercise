FROM alpine:latest

RUN mkdir -p /usr/local/bin

COPY build/_output/cron/expired-order /usr/local/bin/expired-order
RUN chmod +x /usr/local/bin/expired-order

COPY .env /app/.env
RUN sed -i 's/127.0.0.1/host.docker.internal/g' /app/.env

WORKDIR /app

CMD ["sh", "-c", ". /app/.env && /usr/local/bin/expired-order"]
