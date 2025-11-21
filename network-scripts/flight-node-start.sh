#!/bin/sh
set -euo pipefail

RELAY_ENDPOINT=${1:-relay-sidecar-1:8080}
NODE_PRIVATE_KEY=${2:?"missing node private key"}
EVM_RPC_URL=${EVM_RPC_URL:-http://anvil:8545}
FLIGHTS_API_URL=${FLIGHTS_API_URL:-http://flights-api:8085}
FLIGHT_DELAYS_FILE=/deploy-data/flight-delays.address

echo "Waiting for FlightDelays deployment info..."
until [ -s "$FLIGHT_DELAYS_FILE" ]; do sleep 2; done
FLIGHT_DELAYS_ADDRESS=$(cat "$FLIGHT_DELAYS_FILE")

exec /app/flight-node \
  --relay-api-url "http://$RELAY_ENDPOINT" \
  --evm-rpc-url "$EVM_RPC_URL" \
  --flight-delays-address "$FLIGHT_DELAYS_ADDRESS" \
  --flights-api-url "$FLIGHTS_API_URL" \
  --private-key "$NODE_PRIVATE_KEY" \
  --log-level info
