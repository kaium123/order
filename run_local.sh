#!/bin/bash

echo "starting server"

docker compose up consul db cache -d

export ORDERS_CONSUL_PATH="orders"
export ORDERS_CONSUL_URL="localhost:8500"

#export

# Wait for Consul to be ready and upload the configuration
until curl --silent --output /dev/null --fail --request PUT --data-binary @config.docker.yaml http://127.0.0.1:8500/v1/kv/orders; do
    echo "Consul is unavailable or config upload failed - retrying in 2 seconds..."
    sleep 2
done

go run cmd/*.go serve

