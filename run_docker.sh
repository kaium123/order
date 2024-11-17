#!/bin/bash

echo "Starting server..."

# Start dependent services in detached mode
docker compose up consul db cache -d

export ORDERS_CONSUL_PATH="orders"
export ORDERS_CONSUL_URL="localhost:8500"

echo "Uploading configuration to Consul..."

# Wait for Consul to be ready and upload the configuration
until curl --silent --output /dev/null --fail --request PUT --data-binary @config.docker.yaml http://127.0.0.1:8500/v1/kv/orders; do
    echo "Consul is unavailable or config upload failed - retrying in 2 seconds..."
    sleep 2
done

echo "Configuration uploaded successfully to Consul."

# Vendor Go dependencies
echo "Vendoring Go modules..."
go mod vendor

docker rmi -f order-app:latest

docker build -t order-app .


# Clean up vendor directory if needed
 rm -r vendor

# Start the application
echo "Starting the application..."
docker compose up app

