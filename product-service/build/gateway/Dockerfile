FROM alpine:latest

RUN mkdir -p /usr/local/bin

COPY build/_output/gateway /usr/local/bin/gateway
RUN chmod +x /usr/local/bin/gateway

COPY .env /app/.env
RUN sed -i 's/127.0.0.1/host.docker.internal/g' /app/.env

WORKDIR /app

CMD ["sh", "-c", ". /app/.env && /usr/local/bin/gateway"]
