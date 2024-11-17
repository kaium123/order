#!/bin/bash

echo "starting server"

docker compose up consul db cache -d

sleep 10

export ORDERS_CONSUL_PATH="orders"
export ORDERS_CONSUL_URL="localhost:8500"

#export

curl \
    --request PUT \
    --data-binary @config.yaml \
    http://127.0.0.1:8500/v1/kv/orders

go run cmd/*.go serve

