#!/bin/bash

echo "starting server"

export ME_CONSUL_PATH="order"
export ME_CONSUL_URL="localhost:8500"

#export

curl \
    --request PUT \
    --data-binary @config.yaml \
    http://127.0.0.1:8500/v1/kv/order

go run cmd/*.go serve

