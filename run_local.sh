#!/bin/bash

echo "starting server"

export ME_CONSUL_PATH="me"
export ME_CONSUL_URL="localhost:8500"
export GOOGLE_APPLICATION_CREDENTIALS="/$HOME/techetron-service-acc.json"

#export

curl \
    --request PUT \
    --data-binary @config.yaml \
    http://127.0.0.1:8500/v1/kv/me

go run cmd/*.go serve
#go run cmd/*.go listen_account_creation

